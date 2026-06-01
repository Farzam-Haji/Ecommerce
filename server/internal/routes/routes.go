package routes

import (
	"github.com/btalv/Ecommerce/server/internal/handlers"
	"github.com/btalv/Ecommerce/server/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *handlers.ProductHandler) {
	api := r.Group("/api")
	v1 := api.Group("/v1")
	v1.GET("/ping", handlers.Ping)

	v1.GET("/products", handler.AllProducts)
	v1.POST("/products", middlewares.HeaderCheck() ,handler.AddProduct)
	v1.GET("/products/:id", handler.FetchProduct)
	v1.PUT("/products/:id", handler.UpdateProduct)
	v1.PATCH("/products/:id", handler.PatchProduct)
	v1.DELETE("/products/:id", handler.DeleteProduct)
}
