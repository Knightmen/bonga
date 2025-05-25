package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go-server/config"
	_ "go-server/docs" // This is required for swagger
	"go-server/handlers"
	"go-server/models"
)

// @title E-commerce API
// @version 1.0
// @description This is a sample e-commerce server.
// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database
	db := config.ConnectDB()

	// Auto migrate the schema
	db.AutoMigrate(&models.Product{})

	// Initialize router
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "BONGA! JODU!! BONGA!",
		})
	})

	// Initialize product handler
	productHandler := &handlers.ProductHandler{DB: db}

	// Product routes
	v1 := r.Group("/api/v1")
	{
		products := v1.Group("/products")
		{
			products.GET("", productHandler.GetProducts)
			products.GET("/:id", productHandler.GetProduct)
			products.POST("", productHandler.CreateProduct)
			products.PUT("/:id", productHandler.UpdateProduct)
			products.DELETE("/:id", productHandler.DeleteProduct)
		}
	}

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
} 