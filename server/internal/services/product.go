package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/Farzam-Haji/Ecommerce/server/internal/models"
	"github.com/Farzam-Haji/Ecommerce/server/internal/repositories"
)

type ProductRepository interface {
	Create(ctx	context.Context, product *models.Product) (error)
	GetByID(ctx context.Context, id int64) (*models.Product, error)
	GetAll(ctx context.Context) ([]*models.Product, error)
	Update(ctx	context.Context, product *models.Product) (int64, error)
	Delete(ctx	context.Context, id int64, sellerID int64) (int64, error)
	GetBySellerID(ctx context.Context, sellerID int64) ([]*models.Product, error)	
}


type ProductService struct {
	repository  	ProductRepository
}

func NewProductService(repositoryInterface ProductRepository) *ProductService {
	return &ProductService{
		repository: repositoryInterface,
	}
}

func (ps *ProductService) CreateProduct(ctx context.Context, name string, quantity int32, userID int64, price int64) (*models.Product, error) {
	if name == "" {
		return nil, ErrInvalidName
	}
	if quantity < 0 {
		return nil, ErrInvalidQuantity
	}
	if price < 0 {
		return nil, ErrInvalidPrice
	}

	newProduct := models.Product{
		Name: name,
		Quantity: quantity,
		SellerID: userID,
		Price: price,
	}

	err := ps.repository.Create(ctx, &newProduct)
	if errors.Is(err, repositories.ErrDuplicateKey) {
		return nil, ErrDuplicateProduct
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternal, err)
	}

	return &newProduct, nil
}

func (ps *ProductService) GetProduct(ctx context.Context, id int64) (*models.Product, error) {
	if id <=0 {
		return nil, ErrInvalidID
	}

	product, err := ps.repository.GetByID(ctx, id)
	if errors.Is(err, repositories.ErrRowNotFound) {
		return nil, ErrProductNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternal, err)
	}

	return product, nil
}

func (ps *ProductService) GetAllProducts(ctx context.Context) ([]*models.Product, error) {
	products, err := ps.repository.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternal, err)
	}

	return products, nil
}


func (ps *ProductService) UpdateProduct(ctx context.Context, id int64, userID int64, name string, quantity int32, price int64) (*models.Product, error) {
	if id <=0 {
		return nil, ErrInvalidID
	}
	if name == "" {
		return nil, ErrInvalidName
	}
	if quantity < 0 {
		return nil, ErrInvalidQuantity
	}
	if price < 0 {
		return nil, ErrInvalidPrice
	}

	product, err := ps.GetProduct(ctx, id)
	if err != nil {
		return nil, err
	}

	if product.SellerID != userID {
		return nil, ErrForbidden
	}

	product.Name = name
	product.Quantity = quantity
	product.Price = price

	affected, err := ps.repository.Update(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternal, err)
	}
	if affected == 0 {
		return nil, fmt.Errorf("%w: unable to update product %v", ErrInternal, product.Name)
	}

	return product, nil
}

func (ps *ProductService) PatchProduct(ctx context.Context, id int64, userID int64, name *string, quantity *int32, price *int64) (*models.Product, error) {
	if id <= 0 {
		return nil, ErrInvalidID
	}

	product, err := ps.GetProduct(ctx, id)
	if err != nil {
		return nil, err
	}

	if product.SellerID != userID {
		return nil, ErrForbidden
	}

	if name != nil && *name == "" {
		return nil, ErrInvalidName
	}
	if quantity != nil && *quantity < 0 {
		return nil, ErrInvalidQuantity
	}
	if price != nil && *price < 0 {
		return nil, ErrInvalidPrice
	}

	if name != nil {
		product.Name = *name
	}
	if quantity != nil {
		product.Quantity = *quantity
	}
	if price != nil {
		product.Price = *price
	}

	affected, err := ps.repository.Update(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternal, err)
	}
	if affected == 0 {
		return nil, fmt.Errorf("%w: unable to update product %v", ErrInternal, product.Name)
	}

	return product, nil
}

func (ps *ProductService) DeleteProduct(ctx context.Context, id int64, userID int64) (error) {
	if id <= 0 {
		return ErrInvalidID
	}

	product, err := ps.GetProduct(ctx, id)
	if err != nil {
		return err
	}

	if product.SellerID != userID {
		return ErrForbidden
	}

	affected, err := ps.repository.Delete(ctx, id, userID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInternal, err)
	}
	if affected == 0 {
		return fmt.Errorf("%w: unable to delete product %v", ErrInternal, product.Name)
	}

	return nil
}

func (ps *ProductService) GetUserProducts(ctx context.Context, userID int64) ([]*models.Product, error){
	if userID <= 0 {
		return nil, ErrForbidden
	}

	products, err := ps.repository.GetBySellerID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInternal, err)
	}
	
	return products, nil
}