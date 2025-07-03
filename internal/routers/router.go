package routers

import (
	"linkedin-watcher/db"
	"linkedin-watcher/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// RouterDependencies holds all the dependencies needed for routing
type RouterDependencies struct {
	Queries   *db.Queries
	JWTSecret string
}

// SetupRoute creates and configures the main router with proper dependency injection
func SetupRoute(deps *RouterDependencies) *gin.Engine {
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

	// Only register auth routes if database is available
	if deps != nil && deps.Queries != nil {
		RegisterRoutes(router, deps.Queries, deps.JWTSecret)
	} else {
		RegisterRoutesWithoutDB(router)
	}

	return router
}
