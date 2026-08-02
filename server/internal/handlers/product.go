package handlers

import (
	"context"
	"log"
	"strconv"

	"github.com/btalv/Ecommerce/server/internal/models"
	"github.com/btalv/Ecommerce/server/internal/services"
	"github.com/gin-gonic/gin"
)


type ProductServiceAPI interface {
	CreateProduct(context.Context, string, int) (*models.Product, error)
	GetProduct(context.Context, int) (models.Product, error)
	GetAllProducts(context.Context) ([]models.Product, error)
	UpdateProduct(context.Context, int, string, int) 	(error)
	PatchProduct(context.Context, int, *string, *int) (error)
	DeleteProduct(context.Context, int)	(error)	
}

type ProductHandler struct {
	Service 	ProductServiceAPI
}

func NewProductHandler(serviceInterface ProductServiceAPI) *ProductHandler {
	return &ProductHandler{
		Service: serviceInterface,
	}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var request CreateProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println(err)
		c.JSON(400, gin.H{"error": "couldnt bind json"})
		return
	}

	newProduct, err := h.Service.CreateProduct(c.Request.Context(), request.Name, request.Quantity)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(201, gin.H{
		"status":"created",
		"newProduct": *newProduct,
	})

}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	readID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": services.ErrInvalidID.Error()})
		return
	}

	product, err := h.Service.GetProduct(c.Request.Context(), readID)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, gin.H{
		"status": "success",
		"product": product,
	})

}


func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	products, err := h.Service.GetAllProducts(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, gin.H{
		"status": "success",
		"products": products,
	})
}


func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	updateID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": services.ErrInvalidID.Error()})
		return
	}

	var request UpdateProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println(err)
		c.JSON(400, gin.H{"error": "couldnt bind json"})
		return
	}

	err = h.Service.UpdateProduct(c.Request.Context(), updateID, request.Name, request.Quantity)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, gin.H{"status": "success"})
}

func (h *ProductHandler) PatchProduct(c *gin.Context) {
	id := c.Param("id")
	patchID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": services.ErrInvalidID.Error()})
		return
	}

	var request PatchProductRequest
	if err = c.ShouldBindJSON(&request); err != nil {
		log.Println(err)
		c.JSON(400, gin.H{"error": "couldnt bind json"})
		return
	}

	err = h.Service.PatchProduct(c.Request.Context(), patchID, request.Name, request.Quantity)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, gin.H{"status": "success"})
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	deleteID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": services.ErrInvalidID.Error()})
		return
	}

	err = h.Service.DeleteProduct(c.Request.Context(), deleteID)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, gin.H{"status": "success"})
}