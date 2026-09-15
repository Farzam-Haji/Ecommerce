package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/btalv/Ecommerce/server/internal/models"
	"github.com/btalv/Ecommerce/server/internal/repositories"
)

type CartRepository interface {
	CreateCart(ctx context.Context, userID int64) (int64, error)
	GetByUserID(ctx context.Context, userID int64) (*models.Cart, error)
	DeleteCart(ctx context.Context, userID int64) (int64, error)
	AddItem(ctx context.Context, item *models.CartItem) (*models.CartItem, error)
	GetItem(ctx context.Context, cartID int64, productID int64) (*models.CartItem, error)
	GetItemByID(ctx context.Context, itemID int64) (*models.CartItem, error)
	GetCartItems(ctx context.Context, cartID int64) ([]*models.CartItem, error)
	UpdateItemQuantity(ctx context.Context, quantity int32, itemID int64) (int64, error)
	DeleteItem(ctx context.Context, itemID int64) (int64, error)
}

type CartService struct {
	cartRepo 		CartRepository
	productRepo 	ProductRepository
}

func NewCartService(cartRepoInterface CartRepository, productRepoInterface ProductRepository) (*CartService) {
	return &CartService{
		cartRepo: cartRepoInterface,
		productRepo: productRepoInterface,
	}
}

func (cs *CartService) GetUserCart(ctx context.Context, userID int64) (*models.Cart, error) {
	cart, err := cs.cartRepo.GetByUserID(ctx, userID)
	if errors.Is(err, repositories.ErrRowNotFound) {
		return nil, ErrCartNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternal, err)
	}

	return cart, nil
}

func (cs *CartService) DeleteCart(ctx context.Context, userID int64) (error) {
	affected, err := cs.cartRepo.DeleteCart(ctx, userID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInternal, err)
	}
	if affected == 0 {
		return ErrCartNotFound
	}
	return nil
}

func (cs *CartService) DeleteFromCart(ctx context.Context, userID int64, itemID int64) (error) {
	if itemID <= 0 {
		return ErrInvalidID
	}

	item, err := cs.cartRepo.GetItemByID(ctx, itemID)
	if errors.Is(err, repositories.ErrRowNotFound) {
		return ErrItemNotFound
	}
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInternal, err)
	}

	cart, err := cs.GetUserCart(ctx, userID)
	if err != nil {
		return err
	}
	if cart.ID != item.CartID {
		return ErrUnauthorized
	}

	affected, err := cs.cartRepo.DeleteItem(ctx, item.ID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInternal, err)
	}
	if affected == 0 {
		return ErrInternal
	}
	return nil
}

func (cs *CartService) GetUserCartItems(ctx context.Context, userID int64) ([]*models.CartItem, []string, error) {
	cart, err := cs.GetUserCart(ctx, userID)
	if errors.Is(err, repositories.ErrRowNotFound) {
		return nil, nil, ErrCartNotFound
	}
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %w", ErrInternal, err)
	}

	var resultItems []*models.CartItem
	var deletedItems []string
	items, err := cs.cartRepo.GetCartItems(ctx, cart.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %w", ErrInternal, err)
	}

	for _, item := range items {
		if item.ProductID == nil {
			deletedItems = append(deletedItems, item.Name)

			if err := cs.DeleteFromCart(ctx, userID, item.ID); err != nil {
				return nil, nil, err
			}

			continue
		}

		subTotal := int64(item.Quantity) * *item.Price
		item.Price = &subTotal
		resultItems = append(resultItems, item)
	}
	return resultItems, deletedItems, nil
}

func (cs *CartService) AddToCart(ctx context.Context, userID int64, productID int64, pointerQuantity *int32) (*models.CartItem, error) {
	quantity := int32(1)
	if pointerQuantity != nil {
		quantity = *pointerQuantity
	}

	if productID <= 0 {
		return nil, ErrInvalidID
	}
	if quantity < 1 {
		return nil, ErrInvalidQuantity
	}

	product, err := cs.productRepo.GetByID(ctx, productID)
	if errors.Is(err, repositories.ErrRowNotFound) {
		return nil, ErrProductNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternal, err)
	}

	if product.Quantity < quantity {
		return nil, ErrInsufficientQuantity
	}

	var cartID int64
	cart, err := cs.GetUserCart(ctx, userID)
	switch{
	case err == nil:
		cartID = cart.ID

	case errors.Is(err, ErrCartNotFound):
		newID, err := cs.cartRepo.CreateCart(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrInternal, err)
		}

		cartID = newID

	default:
		return nil, err
	}

	item, err := cs.cartRepo.GetItem(ctx, cartID, product.ID)
	if err != nil && errors.Is(err, repositories.ErrRowNotFound) {

		newItem := models.CartItem{
			CartID: cartID,
			ProductID: &product.ID,
			Name: product.Name,
			Quantity: quantity,
		}

		item, err = cs.cartRepo.AddItem(ctx, &newItem)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrInternal, err)
		}

		return item, nil
	}

	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternal, err)
	}


	err = cs.UpdateItemQuantity(ctx, userID, item.ID, item.Quantity+quantity)
	if err != nil {
		return nil, err
	}

	item.Quantity = item.Quantity + quantity
	return item, nil
}

func (cs *CartService) UpdateItemQuantity(ctx context.Context, userID int64, itemID int64, quantity int32) (error) {
	if itemID <= 0 {
		return ErrInvalidID
	}
	if quantity < 0 {
		return ErrInvalidQuantity
	}
	if quantity == 0 {
		return cs.DeleteFromCart(ctx, userID, itemID)
	}

	item, err := cs.cartRepo.GetItemByID(ctx, itemID)
	if errors.Is(err, repositories.ErrRowNotFound) {
		return ErrItemNotFound
	}
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInternal, err)
	}

	cart, err := cs.GetUserCart(ctx, userID)
	if err != nil {
		return err
	}
	
	if cart.ID != item.CartID {
		return ErrUnauthorized
	}

	product, err := cs.productRepo.GetByID(ctx, *item.ProductID)
	if errors.Is(err, repositories.ErrRowNotFound) {
		return ErrProductNotFound
	}
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInternal, err)
	}

	if product.Quantity < quantity {
		return ErrInsufficientQuantity
	}

	affected, err := cs.cartRepo.UpdateItemQuantity(ctx, quantity, item.ID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInternal, err)
	}
	if affected == 0 {
		return ErrInternal
	}
	return nil
}



