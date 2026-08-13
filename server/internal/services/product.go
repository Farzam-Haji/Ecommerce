package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/btalv/Ecommerce/server/internal/models"
)

type ProductRepository interface {
	Create(context.Context, *models.Product) (error)
	GetByID(context.Context, int64) (*models.Product, error)
	GetAll(context.Context) ([]*models.Product, error)
	Update(context.Context, *models.Product) 	(int64, error)
	Delete(context.Context, int64, int64)	(int64, error)
	GetBySellerID(context.Context, int64) ([]*models.Product, error)	
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
	if err != nil {
		return nil, err
	}
	return &newProduct, nil
}

func (ps *ProductService) GetProduct(ctx context.Context, id int64) (*models.Product, error) {
	if id <=0 {
		return nil, ErrInvalidID
	}
	product, err := ps.repository.GetByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrProductNotFound
	}
	return product, err
}

func (ps *ProductService) GetAllProducts(ctx context.Context) ([]*models.Product, error) {
	return ps.repository.GetAll(ctx)
}


func (ps *ProductService) UpdateProduct(ctx context.Context, id int64, userID int64, name string, quantity int32, price int64) (error) {
	if name == "" {
		return ErrInvalidName
	}
	if quantity < 0 {
		return ErrInvalidQuantity
	}

	product, err := ps.GetProduct(ctx, id)
	if err != nil {
		return err
	}
	if product.SellerID != userID {
		return ErrUnauthorized
	}

	product.Name = name
	product.Quantity = quantity
	product.Price = price

	affected, err := ps.repository.Update(ctx, product)
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrInternal
	}
	return nil
}

func (ps *ProductService) PatchProduct(ctx context.Context, id int64, userID int64, name *string, quantity *int32, price *int64) (error) {
	if id <= 0 {
		return ErrInvalidID
	}

	product, err := ps.GetProduct(ctx, id)
	if err != nil {
		return err
	}

	if product.SellerID != userID {
		return ErrUnauthorized
	}

	if name != nil && *name == "" {
		return ErrInvalidName
	}
	if quantity != nil && *quantity < 0 {
		return ErrInvalidQuantity
	}
	if price != nil && *price < 0 {
		return ErrInvalidPrice
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
		return err
	}
	if affected == 0 {
		return ErrInternal
	}
	return nil
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
		return ErrUnauthorized
	}

	affected, err := ps.repository.Delete(ctx, id, userID)
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrInternal
	}
	return nil
}

func (ps *ProductService) GetUserProducts(ctx context.Context, userID int64) ([]*models.Product, error){
	if userID <= 0 {
		return nil, ErrUnauthorized
	}
	return ps.repository.GetBySellerID(ctx, userID)
}