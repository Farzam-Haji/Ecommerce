package routes

import (
	"github.com/Farzam-Haji/Ecommerce/server/internal/web"
	"github.com/gin-gonic/gin"
)

func RegisterWebRoutes(r *gin.Engine, productHandler *web.ProductWebHandler) {

	r.GET("/products/:id", productHandler.GetProduct)
	r.GET("/products", productHandler.GetAllProducts)
	r.GET("/products/:id/preview", productHandler.GetProductPreview)

	r.GET("/tailwind/products", productHandler.GetAllProductsTailwind)
	r.GET("/tailwind/products/:id/preview", productHandler.GetProductPreviewTailwind)
}


	// v1.GET("/products", productHandler.GetAllProducts)
	// v1.POST("/products", middlewares.RequireUser(jwtManager) ,productHandler.CreateProduct)
	// v1.GET("/products/:id", productHandler.GetProduct)
	// v1.PUT("/products/:id", middlewares.RequireUser(jwtManager), productHandler.UpdateProduct)
	// v1.PATCH("/products/:id", middlewares.RequireUser(jwtManager), productHandler.PatchProduct)
	// v1.DELETE("/products/:id", middlewares.RequireUser(jwtManager), productHandler.DeleteProduct)
	// v1.GET("/products/me", middlewares.RequireUser(jwtManager), productHandler.GetMyProducts)