package routers

import (
	"linkedin-watcher/db"
	"linkedin-watcher/internal/controllers"
	"linkedin-watcher/internal/middleware"
	"linkedin-watcher/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RegisterRoutes add all routing list here automatically get main router
func RegisterRoutes(route *gin.Engine, queries *db.Queries, jwtSecret string) {
	// Initialize services
	authService := services.NewAuthService(queries, jwtSecret)

	// Initialize controllers with injected dependencies
	authController := controllers.NewAuthController(authService)

	// Health check endpoint
	route.GET("/health", controllers.HealthCheck)

	// Swagger documentation
	route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Authentication routes
	auth := route.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
		auth.POST("/refresh", authController.RefreshToken)
		auth.POST("/logout", authController.Logout)
	}

	// Protected routes
	protected := route.Group("/auth")
	protected.Use(middleware.AuthMiddleware(authService))
	{
		protected.POST("/change-password", authController.ChangePassword)
	}

	// API v1 routes
	// v1 := route.Group("/api/v1")
	// {
	// 	// Add future API endpoints here
	// 	// v1.GET("/users", handlers.GetUsers)
	// 	// v1.POST("/users", handlers.CreateUser)
	// }

	// 404 handler
	route.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Route Not Found",
			"path":    ctx.Request.URL.Path,
		})
	})
}

// RegisterRoutesWithoutDB registers routes when database is not available
func RegisterRoutesWithoutDB(route *gin.Engine) {
	// Health check endpoint
	route.GET("/health", controllers.HealthCheck)

	// Swagger documentation
	route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 404 handler
	route.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Route Not Found",
			"path":    ctx.Request.URL.Path,
		})
	})
}
