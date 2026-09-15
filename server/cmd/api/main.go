package main

import (
	"fmt"
	"log"
	"os"

	"github.com/btalv/Ecommerce/server/internal/auth"
	"github.com/btalv/Ecommerce/server/internal/database"
	"github.com/btalv/Ecommerce/server/internal/handlers"
	"github.com/btalv/Ecommerce/server/internal/middlewares"
	"github.com/btalv/Ecommerce/server/internal/repositories"
	"github.com/btalv/Ecommerce/server/internal/routes"
	"github.com/btalv/Ecommerce/server/internal/services"
	"github.com/btalv/Ecommerce/server/internal/workers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/net/context"
)


func main() {
    // db initialazion
    err := godotenv.Load()
    if err != nil {
        log.Fatal(err)
    }
    config := database.Config{
        User: os.Getenv("POSTGRES_USER"),
        Password: os.Getenv("POSTGRES_PASSWORD"),
        Host: os.Getenv("POSTGRES_HOST"),
        Port: os.Getenv("POSTGRES_PORT"),
        DBName: os.Getenv("POSTGRES_DBNAME"),
    }

    db, err := database.Open(config)
    if err != nil {
        log.Fatal(err)
    }

    err = db.Ping()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Ping is received. db is connected.")
    defer db.Close()

    // auth constructor
    jwtManager := auth.NewJWTManager(os.Getenv("JWT_SECRET"))

    // product constructor
    postgresProductRepository := repositories.NewPostgresProductRepository(db)
    productService := services.NewProductService(postgresProductRepository)
    productHandler := handlers.NewProductHandler(productService)

    // user constructor
    postgresUserRepository := repositories.NewPostgresUserRepository(db)
    userService := services.NewUserService(postgresUserRepository, jwtManager)
    userHandler := handlers.NewUserHandler(userService)

    // cart constructor
    postgresCartRepository := repositories.NewPostgresCartRepository(db)
    cartService := services.NewCartService(postgresCartRepository, postgresProductRepository)
    cartHandler := handlers.NewCartHandler(cartService)

    // order constructor
    postgresOrderRepository := repositories.NewPostgresOrderRepository(db)
    orderService := services.NewOrderService(postgresOrderRepository)
    orderHandler := handlers.NewOrderHandler(orderService)

    // worker initiation
    workerCTX, workerCancel := context.WithCancel(context.Background())
    defer workerCancel()
    go workers.OrderExpirationWorker(workerCTX, orderService)

    // gin initilaztion
    engine := gin.New()
    engine.Use(gin.Recovery())
    // engine.Use(middlewares.CustomLogger())
    engine.Use(middlewares.RequestID())
    engine.Use(gin.Logger())
    engine.Use(middlewares.ErrorHandler())
    engine.Use(middlewares.Authentication(jwtManager))
    routes.RegisterRoutes(engine, productHandler, userHandler, cartHandler, orderHandler, jwtManager)
    engine.Run() // 0.0.0.0:8080
}