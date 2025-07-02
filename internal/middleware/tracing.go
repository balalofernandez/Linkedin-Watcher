package middleware

import (
	"linkedin-watcher/infra/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// TraceIDKey is the key used to store trace ID in Gin context
	TraceIDKey = "trace_id"
	// TraceIDHeader is the HTTP header name for trace ID
	TraceIDHeader = "X-Trace-ID"
)

// TracingMiddleware adds a unique trace ID to each request
func TracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if trace ID is already provided in headers
		traceID := c.GetHeader(TraceIDHeader)

		// If no trace ID provided, generate a new one
		if traceID == "" {
			traceID = uuid.New().String()
		}

		// Store trace ID in Gin context
		c.Set(TraceIDKey, traceID)

		// Add trace ID to response headers for client tracking
		c.Header(TraceIDHeader, traceID)

		// Log the trace ID for debugging
		logger.Debugf("Request trace ID: %s for path: %s", traceID, c.Request.URL.Path)

		// Continue processing
		c.Next()
	}
}

// GetTraceID retrieves the trace ID from Gin context
func GetTraceID(c *gin.Context) string {
	if traceID, exists := c.Get(TraceIDKey); exists {
		if traceIDStr, ok := traceID.(string); ok {
			return traceIDStr
		}
	}
	return ""
}

// TracingMiddlewareWithCorrelation adds trace ID and correlation ID support
func TracingMiddlewareWithCorrelation() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate or extract trace ID
		traceID := c.GetHeader(TraceIDHeader)
		if traceID == "" {
			traceID = uuid.New().String()
		}

		// Extract correlation ID from headers (for distributed tracing)
		correlationID := c.GetHeader("X-Correlation-ID")
		if correlationID == "" {
			correlationID = uuid.New().String()
		}

		// Store in context
		c.Set(TraceIDKey, traceID)
		c.Set("correlation_id", correlationID)

		// Add to response headers
		c.Header(TraceIDHeader, traceID)
		c.Header("X-Correlation-ID", correlationID)

		// Log for debugging
		logger.Debugf("Request tracing - TraceID: %s, CorrelationID: %s, Path: %s",
			traceID, correlationID, c.Request.URL.Path)

		c.Next()
	}
}
