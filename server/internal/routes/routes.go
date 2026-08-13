package routes

import (
	"github.com/btalv/Ecommerce/server/internal/auth"
	"github.com/btalv/Ecommerce/server/internal/handlers"
	"github.com/btalv/Ecommerce/server/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, productHandler *handlers.ProductHandler, userHandler *handlers.UserHandler, jwtManager *auth.JWTManager) {
	api := r.Group("/api")
	v1 := api.Group("/v1")
	v1.GET("/ping", handlers.Ping)

	v1.GET("/products", productHandler.GetAllProducts)
	v1.POST("/products", middlewares.HeaderCheck(), middlewares.RequireUser(jwtManager) ,productHandler.CreateProduct)
	v1.GET("/products/:id", productHandler.GetProduct)
	v1.PUT("/products/:id", middlewares.RequireUser(jwtManager), productHandler.UpdateProduct)
	v1.PATCH("/products/:id", middlewares.RequireUser(jwtManager), productHandler.PatchProduct)
	v1.DELETE("/products/:id", middlewares.RequireUser(jwtManager), productHandler.DeleteProduct)
	v1.GET("/products/me", middlewares.RequireUser(jwtManager), productHandler.GetMyProducts)

	v1.POST("/register", userHandler.Register)
	v1.POST("/login", userHandler.Login)
}
