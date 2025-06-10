package main

import (
	"log"
	"net/http"

	"go-server/config"
	"go-server/docs"
	"go-server/handlers"
	"go-server/middleware"
	"go-server/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	// Update Swagger host from environment
	docs.SwaggerInfo.Host = config.GetSwaggerHost()

	// Initialize database
	db := config.ConnectDB()

	// Auto migrate the schema
	db.AutoMigrate(&models.Product{}, &models.Resume{})

	// Initialize router
	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "BONGA! JODU!! BONGA!",
		})
	})

	// Initialize handlers
	productHandler := &handlers.ProductHandler{DB: db}
	resumeHandler := handlers.NewResumeHandler(db)
	sessionHandler := handlers.NewSessionHandler(db)

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

		// Resume routes with API key authentication
		resumes := v1.Group("/resume")
		resumes.Use(middleware.APIKeyAuth())
		{
			resumes.GET("/latest", resumeHandler.LatestResume)
			resumes.GET("/getSignedUrl", resumeHandler.GetSignedURL)
			resumes.POST("", resumeHandler.CreateResume)
			resumes.GET("/:id", resumeHandler.GetResume)
			resumes.PUT("/:id", resumeHandler.UpdateResume)
			resumes.DELETE("/:id", resumeHandler.DeleteResume)
		}

		// Session routes with API key authentication
		sessions := v1.Group("/session")
		sessions.Use(middleware.APIKeyAuth())
		{
			sessions.GET("/init", sessionHandler.InitSession)
			sessions.POST("/chat", sessionHandler.ChatSession)
		}
	}

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
} 