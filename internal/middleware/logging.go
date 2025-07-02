package middleware

import (
	"bytes"
	"io"
	"linkedin-watcher/infra/logger"
	"time"

	"github.com/gin-gonic/gin"
)

// responseWriter wraps gin.ResponseWriter to capture response body
type responseWriter struct {
	gin.ResponseWriter
	body   bytes.Buffer
	status int
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *responseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func (w *responseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware logs HTTP requests with detailed information
func LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Extract trace ID from context if available
		traceID := ""
		if traceIDInterface, exists := param.Keys["trace_id"]; exists {
			if traceIDStr, ok := traceIDInterface.(string); ok {
				traceID = traceIDStr
			}
		}

		// Log with structured format
		logger.Infof("HTTP Request - Method: %s, Path: %s, Status: %d, Latency: %s, ClientIP: %s, UserAgent: %s, TraceID: %s, Error: %s",
			param.Method, param.Path, param.StatusCode, param.Latency.String(), param.ClientIP, param.Request.UserAgent(), traceID, param.ErrorMessage)

		return "" // Return empty string as we're using our own logger
	})
}

// RequestLoggingMiddleware provides more detailed request/response logging
func RequestLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Capture request body (for debugging, be careful with sensitive data)
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			// Restore the body for subsequent middleware/handlers
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Wrap response writer to capture response
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			status:         200, // Default status
		}
		c.Writer = writer

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Extract trace ID
		traceID := ""
		if traceIDInterface, exists := c.Get("trace_id"); exists {
			if traceIDStr, ok := traceIDInterface.(string); ok {
				traceID = traceIDStr
			}
		}

		// Log request details
		logFields := map[string]interface{}{
			"method":        c.Request.Method,
			"path":          c.Request.URL.Path,
			"query":         c.Request.URL.RawQuery,
			"status":        writer.status,
			"duration_ms":   duration.Milliseconds(),
			"client_ip":     c.ClientIP(),
			"user_agent":    c.Request.UserAgent(),
			"content_type":  c.GetHeader("Content-Type"),
			"trace_id":      traceID,
			"request_size":  len(requestBody),
			"response_size": writer.body.Len(),
		}

		// Add error information if any
		if len(c.Errors) > 0 {
			logFields["errors"] = c.Errors.String()
		}

		// Log based on status code
		if writer.status >= 400 {
			logger.Errorf("HTTP Request Failed - Method: %s, Path: %s, Status: %d, Duration: %dms, ClientIP: %s, TraceID: %s, Errors: %v",
				logFields["method"], logFields["path"], logFields["status"], logFields["duration_ms"], logFields["client_ip"], logFields["trace_id"], logFields["errors"])
		} else {
			logger.Infof("HTTP Request Completed - Method: %s, Path: %s, Status: %d, Duration: %dms, ClientIP: %s, TraceID: %s",
				logFields["method"], logFields["path"], logFields["status"], logFields["duration_ms"], logFields["client_ip"], logFields["trace_id"])
		}
	}
}
