package handlers

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/btalv/Ecommerce/server/internal/models"
	"github.com/btalv/Ecommerce/server/internal/services"
	"github.com/gin-gonic/gin"
)

var ErrInvalidJSON = errors.New("couldnt bind json")

type ProductServiceAPI interface {
	CreateProduct(context.Context, string, int32, int64, int64) (*models.Product, error)
	GetProduct(context.Context, int64) (*models.Product, error)
	GetAllProducts(context.Context) ([]*models.Product, error)
	UpdateProduct(context.Context, int64, int64, string, int32, int64) 	(error)
	PatchProduct(context.Context, int64, int64, *string, *int32, *int64) (error)
	DeleteProduct(context.Context, int64, int64)	(error)	
	GetUserProducts(context.Context, int64) ([]*models.Product, error)
}

type ProductHandler struct {
	Service 	ProductServiceAPI
}

func NewProductHandler(serviceInterface ProductServiceAPI) *ProductHandler {
	return &ProductHandler{
		Service: serviceInterface,
	}
}

func (ph *ProductHandler) CreateProduct(c *gin.Context) {
	var request CreateProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(fmt.Errorf("%w: %v", ErrInvalidJSON, err))
		return
	}

	value, exists := c.Get("userID")
	if !exists {
		c.Error(services.ErrInternal)
		return
	}
	userID, ok := value.(int64)
	if !ok {
		c.Error(services.ErrInternal)
		return
	}

	newProduct, err := ph.Service.CreateProduct(c.Request.Context(), request.Name, request.Quantity, userID, request.Price)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(201, gin.H{
		"status":"created",
		"newProduct": *newProduct,
	})
}

func (ph *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	productID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": services.ErrInvalidID.Error()})
		return
	}

	product, err := ph.Service.GetProduct(c.Request.Context(), productID)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, gin.H{
		"status": "success",
		"product": product,
	})

}


func (ph *ProductHandler) GetAllProducts(c *gin.Context) {
	products, err := ph.Service.GetAllProducts(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, gin.H{
		"status": "success",
		"products": products,
	})
}


func (ph *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	productID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": services.ErrInvalidID.Error()})
		return
	}

	var request UpdateProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(fmt.Errorf("%w: %v", ErrInvalidJSON, err))
		return
	}

	value, exists := c.Get("userID")
	if !exists {
		c.Error(services.ErrInternal)
		return
	}
	userID, ok := value.(int64)
	if !ok {
		c.Error(services.ErrInternal)
		return
	}

	err = ph.Service.UpdateProduct(c.Request.Context(), productID, userID, request.Name, request.Quantity, request.Price)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, gin.H{"status": "success"})
}

func (ph *ProductHandler) PatchProduct(c *gin.Context) {
	id := c.Param("id")
	productID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": services.ErrInvalidID.Error()})
		return
	}

	var request PatchProductRequest
	if err = c.ShouldBindJSON(&request); err != nil {
		c.Error(fmt.Errorf("%w: %v", ErrInvalidJSON, err))
		return
	}

	value, exists := c.Get("userID")
	if !exists {
		c.Error(services.ErrInternal)
		return
	}
	userID, ok := value.(int64)
	if !ok {
		c.Error(services.ErrInternal)
		return
	}

	err = ph.Service.PatchProduct(c.Request.Context(), productID, userID, request.Name, request.Quantity, request.Price)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, gin.H{"status": "success"})
}

func (ph *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	productID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": services.ErrInvalidID.Error()})
		return
	}

	value, exists := c.Get("userID")
	if !exists {
		c.Error(services.ErrInternal)
		return
	}
	userID, ok := value.(int64)
	if !ok {
		c.Error(services.ErrInternal)
		return
	}

	err = ph.Service.DeleteProduct(c.Request.Context(), productID, userID)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, gin.H{"status": "success"})
}

func (ph *ProductHandler) GetMyProducts(c *gin.Context) {
	value, exists := c.Get("userID")
	if !exists {
		c.Error(services.ErrInternal)
		return
	}
	userID, ok := value.(int64)
	if !ok {
		c.Error(services.ErrInternal)
		return
	}

	products, err := ph.Service.GetUserProducts(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, gin.H{
		"status":"success",
		"products":products,
	})
}