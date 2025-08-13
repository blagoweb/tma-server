package routes

import (
	"github.com/gin-gonic/gin"
	"tma/handlers"
	"tma/middleware"
	"tma/auth"
)

func SetupRoutes(
	authHandler *handlers.AuthHandler,
	itemsHandler *handlers.ItemsHandler,
	jwtManager *auth.JWTManager,
) *gin.Engine {
	router := gin.Default()

	// Add CORS middleware
	router.Use(middleware.CORSMiddleware())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"message": "Telegram Mini App API is running",
		})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Public routes (no authentication required)
		auth := api.Group("/auth")
		{
			auth.POST("/telegram", authHandler.Auth)
			auth.POST("/telegram/test", authHandler.TestAuth) // Test endpoint without hash validation
		}

		// Protected routes (authentication required)
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(jwtManager))
		{
			// User routes
			user := protected.Group("/user")
			{
				user.GET("/profile", authHandler.GetProfile)
			}

			// Items routes
			items := protected.Group("/items")
			{
				items.GET("", itemsHandler.GetItems)
				items.GET("/:id", itemsHandler.GetItem)
				items.POST("", itemsHandler.CreateItem)
				items.PUT("/:id", itemsHandler.UpdateItem)
				items.DELETE("/:id", itemsHandler.DeleteItem)
			}
		}
	}

	return router
}
