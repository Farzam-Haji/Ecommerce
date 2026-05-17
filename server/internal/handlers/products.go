package handlers

import (
	"log"
	"strconv"
	"github.com/btalv/Ecommerce/server/internal/models"
	"github.com/gin-gonic/gin"
)


type ProductRepository interface {
	Create(models.Product) (int, error)
	Read(int) (models.Product, error)
	ReadAll() ([]models.Product, error)
	Update(models.Product) 	(int64, error)
	Delete(int)	(int64, error)	
}

type ProductHandler struct {
	Repository 	ProductRepository
}

func NewProductHandler(repo ProductRepository) *ProductHandler {
	return &ProductHandler{
		Repository: repo,
	}
}

func (h *ProductHandler) AllProducts(c *gin.Context) {
	products, err := h.Repository.ReadAll()
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": "something went wrong"})
		return
	}
	c.JSON(200, gin.H{
		"status": "success",
		"products": products,
	})
}

func (h *ProductHandler) FetchProduct(c *gin.Context) {
	id := c.Param("id")
	readID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "wrong id format",
		})
		return
	}

	product, err := h.Repository.Read(readID)
	if err != nil {
		c.JSON(404, gin.H{"error": "something went wrong while fetching the product"})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"product": product,
	})
}

func (h *ProductHandler) AddProduct(c *gin.Context) {
	var request CreateProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println(err)
		c.JSON(400, gin.H{"error": "couldn't bind json"})
		return
	}
	newProduct := models.Product{
		Name: request.Name,
		Quantity: request.Quantity,
	}

	newID, err := h.Repository.Create(newProduct)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": "couldn't add to database"})
		return
	}

	c.JSON(201, gin.H{
		"status": "created",
		"newProductID": newID,
	})
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	updateID, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{"error": "wrong id format"})
		return
	}

	var request UpdateProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println(err)
		c.JSON(400, gin.H{"error": "couldn't bind json"})
		return
	}
	updatingProduct := models.Product{
		ID: updateID,
		Name: *request.Name,
		Quantity: *request.Quantity,
	}

	affected, err := h.Repository.Update(updatingProduct)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": "couldn't update database"})
		return
	}
	if affected == 0 {
		c.JSON(404, gin.H{"error": "id don't exist"})
		return
	}
	c.JSON(200, gin.H{"status": "success"})
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	deleteID, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{"error": "wrong id format"})
		return
	}

	affected, err := h.Repository.Delete(deleteID)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": "couldn't update database"})
		return
	}
	if affected == 0 {
		c.JSON(404, gin.H{"error": "couldn't find the data"})
		return
	}
	c.JSON(200, gin.H{"status": "success"})
}