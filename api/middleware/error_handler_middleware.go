package middleware

import (
	"errors"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/lyj404/gin-api-template/domain/result"
	"go.uber.org/zap"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	TraceID string `json:"trace_id,omitempty"`
}

func ErrorHandlerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			lastErr := c.Errors.Last()

			var statusCode int
			var message string

			if lastErr.Type == gin.ErrorTypeBind {
				statusCode = http.StatusBadRequest
				message = "Invalid request parameters"
			} else if lastErr.Type == gin.ErrorTypePrivate {
				statusCode = http.StatusInternalServerError
				message = "Internal server error"
			} else {
				statusCode = http.StatusInternalServerError
				message = lastErr.Error()
			}

			logger.Error("Request error",
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
				zap.String("error", lastErr.Error()),
				zap.String("trace_id", c.GetString("trace_id")),
			)

			if !c.Writer.Written() {
				result.ErrorResponse(c, statusCode, message)
			}
		}
	}
}

func RecoveryMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()
				traceID := c.GetString("trace_id")

				logger.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
					zap.String("trace_id", traceID),
					zap.String("stack", string(stack)),
				)

				if !c.Writer.Written() {
					result.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
				}

				c.Abort()
			}
		}()
		c.Next()
	}
}

func NewError(err error) error {
	return errors.New(err.Error())
}
