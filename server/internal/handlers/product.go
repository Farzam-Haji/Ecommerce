package handlers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Farzam-Haji/Ecommerce/server/internal/models"
	"github.com/gin-gonic/gin"
)


type ProductServiceAPI interface {
	CreateProduct(context.Context, string, int32, int64, int64) (*models.Product, error)
	GetProduct(context.Context, int64) (*models.Product, error)
	GetAllProducts(context.Context) ([]*models.Product, error)
	UpdateProduct(context.Context, int64, int64, string, int32, int64) 	(*models.Product, error)
	PatchProduct(context.Context, int64, int64, *string, *int32, *int64) (*models.Product, error)
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
		c.Error(fmt.Errorf("%w: %w", ErrInvalidJSON, err))
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	newProduct, err := ph.Service.CreateProduct(c.Request.Context(), request.Name, *request.Quantity, userID, *request.Price)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(201, gin.H{
		"product": newProduct,
	})
}

func (ph *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	productID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.Error(fmt.Errorf("%w: %w", ErrInvalidParam, err))
		return
	}

	product, err := ph.Service.GetProduct(c.Request.Context(), productID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{
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
		"products": products,
	})
}


func (ph *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	productID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.Error(fmt.Errorf("%w: %w", ErrInvalidParam, err))
		return
	}

	var request UpdateProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(fmt.Errorf("%w: %w", ErrInvalidJSON, err))
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	updatedProduct, err := ph.Service.UpdateProduct(c.Request.Context(), productID, userID, request.Name, *request.Quantity, *request.Price)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{
		"product": updatedProduct,
	})
}

func (ph *ProductHandler) PatchProduct(c *gin.Context) {
	id := c.Param("id")
	productID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.Error(fmt.Errorf("%w: %w", ErrInvalidParam, err))
		return
	}

	var request PatchProductRequest
	if err = c.ShouldBindJSON(&request); err != nil {
		c.Error(fmt.Errorf("%w: %w", ErrInvalidJSON, err))
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	updatedProduct, err := ph.Service.PatchProduct(c.Request.Context(), productID, userID, request.Name, request.Quantity, request.Price)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{
		"product": updatedProduct,
	})
}

func (ph *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	productID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.Error(fmt.Errorf("%w: %w", ErrInvalidParam, err))
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	err = ph.Service.DeleteProduct(c.Request.Context(), productID, userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(204)
}

func (ph *ProductHandler) GetMyProducts(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	products, err := ph.Service.GetUserProducts(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{
		"products":products,
	})
}