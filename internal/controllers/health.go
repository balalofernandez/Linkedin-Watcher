package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string    `json:"status" example:"ok"`
	Timestamp time.Time `json:"timestamp" example:"2024-01-01T00:00:00Z"`
	Service   string    `json:"service" example:"linkedin-watcher"`
	Version   string    `json:"version" example:"1.0.0"`
	Uptime    string    `json:"uptime" example:"1h30m45s"`
}

// HealthCheck godoc
// @Summary Health check endpoint
// @Description Returns the health status of the application
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	// In a real application, you might want to check:
	// - Database connectivity
	// - External service dependencies
	// - Application metrics

	response := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now(),
		Service:   "linkedin-watcher",
		Version:   "1.0.0",
		Uptime:    "1h30m45s", // This would be calculated in a real app
	}

	c.JSON(http.StatusOK, response)
}
