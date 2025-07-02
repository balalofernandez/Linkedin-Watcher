package routers

import (
	"linkedin-watcher/internal/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RegisterRoutes add all routing list here automatically get main router
func RegisterRoutes(route *gin.Engine) {
	// Health check endpoint
	route.GET("/health", controllers.HealthCheck)

	// Swagger documentation
	route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
