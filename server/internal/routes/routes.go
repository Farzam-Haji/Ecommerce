package routes

import (
	"github.com/btalv/Ecommerce/server/internal/auth"
	"github.com/btalv/Ecommerce/server/internal/handlers"
	"github.com/btalv/Ecommerce/server/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, productHandler *handlers.ProductHandler, userHandler *handlers.UserHandler, cartHandler *handlers.CartHandler, orderHandler *handlers.OrderHandler, jwtManager *auth.JWTManager) {
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

	v1.GET("/cart", middlewares.RequireUser(jwtManager), cartHandler.GetMyCartItems)
	v1.DELETE("/cart", middlewares.RequireUser(jwtManager), cartHandler.DeleteMyCart)
	v1.POST("/cart/items", middlewares.RequireUser(jwtManager), cartHandler.AddItemToMyCart)
	v1.PATCH("/cart/items/:id", middlewares.RequireUser(jwtManager), cartHandler.UpdateMyItemQuantity)
	v1.DELETE("/cart/items/:id", middlewares.RequireUser(jwtManager), cartHandler.DeleteMyItem)

	v1.GET("/orders", middlewares.RequireUser(jwtManager), orderHandler.GetMyOrders)
	v1.POST("/orders", middlewares.RequireUser(jwtManager), orderHandler.CreateOrder)
	v1.GET("/orders/:id", middlewares.RequireUser(jwtManager), orderHandler.GetOrder)
}
