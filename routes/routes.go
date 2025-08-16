package routes

import (
	"tma/auth"
	"tma/handlers"
	"tma/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	authHandler *handlers.AuthHandler,
	pagesHandler *handlers.PagesHandler,
	jwtManager *auth.JWTManager,
) *gin.Engine {
	router := gin.Default()

	// Add CORS middleware
	router.Use(middleware.CORSMiddleware())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
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

		// Public page route (no authentication required)
		api.GET("/pages/:id", pagesHandler.GetPage)

		// Protected routes (authentication required)
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(jwtManager))
		{
			// User routes
			user := protected.Group("/user")
			{
				user.GET("/profile", authHandler.GetProfile)
			}

			// Protected pages routes
			protectedPages := protected.Group("/pages")
			{
				protectedPages.GET("", pagesHandler.GetPages)
				protectedPages.POST("", pagesHandler.CreatePage)
				protectedPages.PUT("/:id", pagesHandler.UpdatePage)
				protectedPages.DELETE("/:id", pagesHandler.DeletePage)
			}
		}
	}

	return router
}
