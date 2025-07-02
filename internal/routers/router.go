package routers

import (
	"linkedin-watcher/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func SetupRoute() *gin.Engine {

	environment := viper.GetBool("DEBUG")
	if environment {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	allowedHosts := viper.GetString("ALLOWED_HOSTS")
	router := gin.New()
	router.SetTrustedProxies([]string{allowedHosts})

	// Add middleware in order of execution
	router.Use(middleware.TracingMiddleware())        // First: Add trace ID
	router.Use(middleware.RequestLoggingMiddleware()) // Second: Log requests with trace ID
	router.Use(gin.Recovery())                        // Third: Recovery middleware
	router.Use(middleware.CORSMiddleware())           // Fourth: CORS middleware

	RegisterRoutes(router) //routes register

	return router
}
