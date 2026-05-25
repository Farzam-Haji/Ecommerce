package services

import (
	"database/sql"
	"errors"
	"github.com/btalv/Ecommerce/server/internal/models"
)

type ProductRepository interface {
	Create(models.Product) (int, error)
	Read(int) (models.Product, error)
	ReadAll() ([]models.Product, error)
	Update(models.Product) 	(int64, error)
	Delete(int)	(int64, error)	
}


type productService struct {
	repository  	ProductRepository
}

func NewProductService(repositoryInterface ProductRepository) *productService {
	return &productService{
		repository: repositoryInterface,
	}
}

func (s *productService) CreateProduct(request models.Product) (int, error) {
	if request.Name == "" {
		return 0, ErrInvalidName
	}
	if request.Quantity < 0 {
		return 0, ErrInvalidQuantity
	}
	return s.repository.Create(request)
}

func (s *productService) GetProduct(request int) (models.Product, error) {
	if request <=0 {
		return models.Product{}, ErrInvalidID
	}
	product, err := s.repository.Read(request)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Product{}, ErrProductNotFound
	}
	return product, err
}

func (s *productService) GetProducts() ([]models.Product, error) {
	return s.repository.ReadAll()
}


func (s *productService) UpdateProduct(request models.Product) (error) {
	if request.Name == "" {
		return ErrInvalidName
	}
	if request.Quantity < 0 {
		return ErrInvalidQuantity
	}

	affected, err := s.repository.Update(request)
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrProductNotFound
	}
	return nil
}

func (s *productService) PatchProduct(id int, name *string, quantity *int) (error) {
	if id <= 0 {
		return ErrInvalidID
	}

	updatingProduct, err := s.repository.Read(id)
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

	affected, err := s.repository.Update(updatingProduct)
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrProductNotFound
	}
	return nil
}

func (s *productService) DeleteProduct(request int) (error) {
	if request <= 0 {
		return ErrInvalidID
	}

	affected, err := s.repository.Delete(request)
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrProductNotFound
	}
	return nil
}