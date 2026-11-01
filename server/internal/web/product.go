package web

import (
	"context"
	// "fmt"
	"strconv"

	"github.com/Farzam-Haji/Ecommerce/server/internal/models"
	"github.com/gin-gonic/gin"
)

type ProductServiceWeb interface {
	GetProduct(ctx context.Context, id int64) (*models.Product, error)
	GetAllProducts(ctx context.Context) ([]*models.Product, error)
}

type ProductWebHandler struct {
	service 	ProductServiceWeb
}

func NewProductWebHandler(serviceInterface ProductServiceWeb) *ProductWebHandler {
	return &ProductWebHandler{
		service: serviceInterface,
	}
}

func (ph *ProductWebHandler) GetProduct(c *gin.Context) {
	param := c.Param("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		c.JSON(400, "bad param")
		return
	}

	product, err := ph.service.GetProduct(c.Request.Context(), id)
	if err != nil {
		c.HTML(500, "error.html", err)
		return
	}

	c.HTML(200, "product.html", product)
}

func (ph *ProductWebHandler) GetAllProducts(c *gin.Context) {
	products, err := ph.service.GetAllProducts(c.Request.Context())
	if err != nil {
		c.HTML(500, "error.html", err)
		return
	}

	c.HTML(200, "products.html", products)
}

func (ph *ProductWebHandler) GetProductPreview(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, "bad param")
		return
	}

	product, err := ph.service.GetProduct(c.Request.Context(), id)
	if err != nil {
		c.HTML(500, "error.html", err)
		return
	}

	c.HTML(200, "productPreview.html", product)
}

func (ph *ProductWebHandler) GetAllProductsTailwind(c *gin.Context) {
    products, err := ph.service.GetAllProducts(c.Request.Context())
    if err != nil {
        c.HTML(500, "error.html", gin.H{
            "error": err.Error(),
        })
        return
    }

    c.HTML(200, "products-tailwind.html", products)
}

func (ph *ProductWebHandler) GetProductPreviewTailwind(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, "bad param")
		return
	}

	product, err := ph.service.GetProduct(c.Request.Context(), id)
	if err != nil {
		c.HTML(500, "error.html", err)
		return
	}

	c.HTML(200, "productPreview-tailwind.html", product)
}