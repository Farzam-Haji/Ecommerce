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
	GetCartItem(ctx context.Context, cartID int64, productID int64) (*models.CartItem, error)
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
		return ErrForbidden
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
	resultItems := make([]*models.CartItem, 0)

	cart, err := cs.GetUserCart(ctx, userID)
	if errors.Is(err, ErrCartNotFound) {
		return resultItems, nil, ErrCartNotFound
	}
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %w", ErrInternal, err)
	}

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

	item, err := cs.cartRepo.GetCartItem(ctx, cartID, product.ID)
	if err != nil && errors.Is(err, repositories.ErrRowNotFound) {
		subTotal := int64(quantity) * product.Price

		newItem := models.CartItem{
			CartID: cartID,
			ProductID: &product.ID,
			Name: product.Name,
			Quantity: quantity,
			Price: &subTotal,
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


	item, err = cs.UpdateItemQuantity(ctx, userID, item.ID, item.Quantity+quantity)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (cs *CartService) UpdateItemQuantity(ctx context.Context, userID int64, itemID int64, quantity int32) (*models.CartItem, error) {
	if itemID <= 0 {
		return nil, ErrInvalidID
	}
	if quantity < 0 {
		return nil, ErrInvalidQuantity
	}
	if quantity == 0 {
		return nil, cs.DeleteFromCart(ctx, userID, itemID)
	}

	item, err := cs.cartRepo.GetItemByID(ctx, itemID)
	if errors.Is(err, repositories.ErrRowNotFound) {
		return nil, ErrItemNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternal, err)
	}

	cart, err := cs.GetUserCart(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	if cart.ID != item.CartID {
		return nil, ErrForbidden
	}

	product, err := cs.productRepo.GetByID(ctx, *item.ProductID)
	if errors.Is(err, repositories.ErrRowNotFound) {
		return nil, ErrProductNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternal, err)
	}

	if product.Quantity < quantity {
		return nil, ErrInsufficientQuantity
	}

	affected, err := cs.cartRepo.UpdateItemQuantity(ctx, quantity, item.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternal, err)
	}
	if affected == 0 {
		return nil, ErrInternal
	}

	item.Quantity = quantity
	subTotal := int64(item.Quantity) * product.Price
	item.Price = &subTotal
	return item, nil
}



