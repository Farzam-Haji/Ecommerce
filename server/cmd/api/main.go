package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/Farzam-Haji/Ecommerce/server/internal/auth"
	"github.com/Farzam-Haji/Ecommerce/server/internal/bank"
	"github.com/Farzam-Haji/Ecommerce/server/internal/database"
	"github.com/Farzam-Haji/Ecommerce/server/internal/handlers"
	"github.com/Farzam-Haji/Ecommerce/server/internal/middlewares"
	"github.com/Farzam-Haji/Ecommerce/server/internal/repositories"
	"github.com/Farzam-Haji/Ecommerce/server/internal/routes"
	"github.com/Farzam-Haji/Ecommerce/server/internal/services"
	"github.com/Farzam-Haji/Ecommerce/server/internal/web"
	"github.com/Farzam-Haji/Ecommerce/server/internal/workers"
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

    // bank client constructor
    client := http.Client{}
    bankClient := bank.NewBankClient(os.Getenv("BANK_URL"), &client)

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
    orderService := services.NewOrderService(postgresOrderRepository, bankClient)
    orderHandler := handlers.NewOrderHandler(orderService, os.Getenv("SERVER"))

    // worker initiation
    workerCTX, workerCancel := context.WithCancel(context.Background())
    defer workerCancel()
    go workers.OrderExpirationWorker(workerCTX, orderService)

    // web init
    productWebHandler := web.NewProductWebHandler(productService) 

    // gin initilaztion
    engine := gin.New()

    engine.Use(gin.Recovery())
    engine.Use(middlewares.RequestID())
    engine.Use(gin.Logger())
    engine.Use(middlewares.ErrorHandler())

    engine.Static("/static", "./static")

    engine.SetFuncMap(template.FuncMap{
        "formatPrice": func(price int64) string {
            return fmt.Sprintf("$ %.2f", float64(price)/100)
        },
    })
        
    engine.LoadHTMLGlob("templates/*")

    routes.RegisterRoutes(engine, productHandler, userHandler, cartHandler, orderHandler, jwtManager)
    routes.RegisterWebRoutes(engine, productWebHandler)
    engine.Run(os.Getenv("SERVER")) // 0.0.0.0:8080
}