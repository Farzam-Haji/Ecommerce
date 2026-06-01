package main

import (
	"fmt"
    "os"
    "log"
	"github.com/btalv/Ecommerce/server/internal/database"
	"github.com/btalv/Ecommerce/server/internal/routes"
    "github.com/btalv/Ecommerce/server/internal/handlers"
    "github.com/btalv/Ecommerce/server/internal/repositories"
    "github.com/btalv/Ecommerce/server/internal/services"
    "github.com/btalv/Ecommerce/server/internal/middlewares"
	"github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
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

    //constructor
    postgresProductRepository := repositories.NewPostgresProductRepository(db)
    productService := services.NewProductService(postgresProductRepository)
    productHandler := handlers.NewProductHandler(productService)

    // gin initilaztion
    engine := gin.New()
    engine.Use(middlewares.CustomLogger())
    engine.Use(middlewares.RequestID())
    engine.Use(gin.Logger())
    engine.Use(gin.Recovery())
    engine.Use(middlewares.ErrorHandler())
    routes.RegisterRoutes(engine, productHandler)
    engine.Run() // 0.0.0.0:8080
}