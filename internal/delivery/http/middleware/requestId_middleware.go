package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const requestIDKey = "requestId"
const requestIDHeader = "X-Request-ID"

func RequestIDMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reqID := ctx.GetHeader(requestIDHeader)
		if reqID == "" {
			reqID = uuid.NewString()
		}

		ctx.Set(requestIDKey, reqID)
		ctx.Header(requestIDHeader, reqID)

		ctx.Next()
	}
}
