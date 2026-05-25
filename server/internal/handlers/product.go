package handlers

import (
	"errors"
	"log"
	"strconv"

	"github.com/btalv/Ecommerce/server/internal/models"
	"github.com/btalv/Ecommerce/server/internal/services"
	"github.com/gin-gonic/gin"
)


type ProductService interface {
	CreateProduct(models.Product) (int, error)
	GetProduct(int) (models.Product, error)
	GetProducts() ([]models.Product, error)
	UpdateProduct(models.Product) 	(error)
	DeleteProduct(int)	(error)	
	PatchProduct(int, *string, *int) (error)
}

type ProductHandler struct {
	Service 	ProductService
}

func NewProductHandler(serviceInterface ProductService) *ProductHandler {
	return &ProductHandler{
		Service: serviceInterface,
	}
}

func (h *ProductHandler) AddProduct(c *gin.Context) {
	var request CreateProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println(err)
		c.JSON(400, gin.H{"error": "couldnt bind json"})
		return
	}
	newProduct := models.Product{
		Name: request.Name,
		Quantity: request.Quantity,
	}

	newID, err := h.Service.CreateProduct(newProduct)
	switch {
	case err == nil:
		c.JSON(201, gin.H{
			"status":"created",
			"newProductID": newID,
		})
	case errors.Is(err, services.ErrInvalidName):
		c.JSON(400, gin.H{"error":err.Error()})
	case errors.Is(err, services.ErrInvalidQuantity):
		c.JSON(400, gin.H{"error":err.Error()})
	default:
		log.Println(err)
		c.JSON(500, gin.H{"error":services.ErrDatabase.Error()})
	}
}

func (h *ProductHandler) FetchProduct(c *gin.Context) {
	id := c.Param("id")
	readID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": services.ErrInvalidID.Error()})
		return
	}

	product, err := h.Service.GetProduct(readID)
	switch{
	case err == nil:
		c.JSON(200, gin.H{
		"status": "success",
		"product": product,
		})
	case errors.Is(err, services.ErrInvalidID):
		c.JSON(400, gin.H{"error": services.ErrInvalidID.Error()})
	case errors.Is(err, services.ErrProductNotFound):
		c.JSON(404, gin.H{"error": services.ErrProductNotFound.Error()})
	default:
		log.Println(err)
		c.JSON(500, gin.H{"error": services.ErrDatabase.Error()})
	}
}


func (h *ProductHandler) AllProducts(c *gin.Context) {
	products, err := h.Service.GetProducts()
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": services.ErrDatabase.Error()})
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
	updatingProduct := models.Product{
		ID: updateID,
		Name: request.Name,
		Quantity: request.Quantity,
	}

	err = h.Service.UpdateProduct(updatingProduct)
	switch{
	case err == nil:
		c.JSON(200, gin.H{"status": "success"})
	case errors.Is(err, services.ErrInvalidName):
		c.JSON(400, gin.H{"error": services.ErrInvalidName.Error()})
	case errors.Is(err, services.ErrInvalidQuantity):
		c.JSON(400, gin.H{"error": services.ErrInvalidQuantity.Error()})
	case errors.Is(err, services.ErrProductNotFound):
		c.JSON(404, gin.H{"error": services.ErrProductNotFound.Error()})
	default:
		log.Println(err)
		c.JSON(500, gin.H{"error": services.ErrDatabase.Error()})
	}
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

	err = h.Service.PatchProduct(patchID, request.Name, request.Quantity)
	switch{
	case err == nil:
		c.JSON(200, gin.H{"status": "success"})
	case errors.Is(err, services.ErrInvalidID):
		c.JSON(400, gin.H{"error": services.ErrInvalidID.Error()})
	case errors.Is(err, services.ErrProductNotFound):
		c.JSON(404, gin.H{"error": services.ErrProductNotFound.Error()})
	case errors.Is(err, services.ErrInvalidName):
		c.JSON(400, gin.H{"error": services.ErrInvalidName.Error()})
	case errors.Is(err, services.ErrInvalidQuantity):
		c.JSON(400, gin.H{"error": services.ErrInvalidQuantity.Error()})
	default:
		log.Println(err)
		c.JSON(500, gin.H{"error": services.ErrDatabase.Error()})
	}
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	deleteID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": services.ErrInvalidID.Error()})
		return
	}

	err = h.Service.DeleteProduct(deleteID)
	switch{
	case err == nil:
		c.JSON(200, gin.H{"status": "success"})
	case errors.Is(err, services.ErrProductNotFound):
		c.JSON(404, gin.H{"error": services.ErrProductNotFound.Error()})
	default:
		log.Println(err)
		c.JSON(500, gin.H{"error": services.ErrDatabase.Error()})
	}
}