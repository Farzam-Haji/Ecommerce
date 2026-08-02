package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/btalv/Ecommerce/server/internal/models"
)

type ProductRepository interface {
	Create(context.Context, *models.Product) (error)
	GetByID(context.Context, int) (models.Product, error)
	GetAll(context.Context) ([]models.Product, error)
	Update(context.Context, models.Product) 	(int64, error)
	Delete(context.Context, int)	(int64, error)	
}


type ProductService struct {
	repository  	ProductRepository
}

func NewProductService(repositoryInterface ProductRepository) *ProductService {
	return &ProductService{
		repository: repositoryInterface,
	}
}

func (s *ProductService) CreateProduct(ctx context.Context, name string, quantity int) (*models.Product, error) {
	if name == "" {
		return nil, ErrInvalidName
	}
	if quantity < 0 {
		return nil, ErrInvalidQuantity
	}
	newProduct := models.Product{
		Name: name,
		Quantity: quantity,
	}

	err := s.repository.Create(ctx, &newProduct)
	if err != nil {
		return nil, err
	}
	return &newProduct, nil
}

func (s *ProductService) GetProduct(ctx context.Context, request int) (models.Product, error) {
	if request <=0 {
		return models.Product{}, ErrInvalidID
	}
	product, err := s.repository.GetByID(ctx, request)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Product{}, ErrProductNotFound
	}
	return product, err
}

func (s *ProductService) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	return s.repository.GetAll(ctx)
}


func (s *ProductService) UpdateProduct(ctx context.Context, id int, name string, quantity int) (error) {
	if name == "" {
		return ErrInvalidName
	}
	if quantity < 0 {
		return ErrInvalidQuantity
	}

	product := models.Product{
		ID: id,
		Name: name,
		Quantity: quantity,
	}

	affected, err := s.repository.Update(ctx, product)
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrProductNotFound
	}
	return nil
}

func (s *ProductService) PatchProduct(ctx context.Context, id int, name *string, quantity *int) (error) {
	if id <= 0 {
		return ErrInvalidID
	}

	updatingProduct, err := s.repository.GetByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrProductNotFound
	}
	if err != nil {
		return err
	}

	if name != nil && *name == "" {
		return ErrInvalidName
	}
	if quantity != nil && *quantity < 0 {
		return ErrInvalidQuantity
	}

	if name != nil {
		updatingProduct.Name = *name
	}
	if quantity != nil {
		updatingProduct.Quantity = *quantity
	}

	affected, err := s.repository.Update(ctx, updatingProduct)
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrProductNotFound
	}
	return nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, request int) (error) {
	if request <= 0 {
		return ErrInvalidID
	}

	affected, err := s.repository.Delete(ctx, request)
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrProductNotFound
	}
	return nil
}