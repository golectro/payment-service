package middleware

import (
	"context"
	"golectro-payment/internal/usecase"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoggingMiddleware(logger *logrus.Logger, logUC *usecase.LogUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		reqID, _ := c.Get("requestId")

		fields := logrus.Fields{
			"status":    c.Writer.Status(),
			"method":    c.Request.Method,
			"path":      path,
			"query":     query,
			"ip":        c.ClientIP(),
			"userAgent": c.Request.UserAgent(),
			"latency":   latency,
			"requestId": reqID,
		}

		if gin.Mode() == gin.DebugMode {
			fields["headers"] = c.Request.Header
		}

		auth := GetUser(c)
		var userID string
		if auth != nil {
			userID = auth.ID.String()
		}

		message := "Request to " + path

		if userID != "" {
			if reqIDStr, ok := reqID.(string); ok {
				go func(ctx context.Context, userID, path, reqID string, status int, errStr string) {
					level := "INFO"
					if status >= 500 {
						level = "ERROR"
					} else if status >= 400 {
						level = "WARN"
					}
					if err := logUC.LogActivity(context.Background(), level, reqID, message, userID, path, status, errStr); err != nil {
						logger.WithError(err).Warn("Failed to log activity to MongoDB")
					}
				}(c.Request.Context(), userID, path, reqIDStr, c.Writer.Status(), c.Errors.ByType(gin.ErrorTypePrivate).String())
			}
		}

		entry := logger.WithFields(fields)
		switch {
		case c.Writer.Status() >= 500:
			entry.Error(message)
		case c.Writer.Status() >= 400:
			entry.Warn(message)
		default:
			entry.Info(message)
		}
	}
}
