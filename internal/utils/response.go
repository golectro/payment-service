package utils

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"golectro-payment/internal/model"

	"github.com/gin-gonic/gin"
)

func isMultilangError(errMsg string) bool {
	return strings.Contains(errMsg, ":") && strings.Contains(errMsg, "|")
}

func extractErrorMessage(fullErr string) string {
	if idx := strings.Index(fullErr, ": "); idx != -1 {
		return strings.TrimSpace(fullErr[idx+2:])
	}
	return fullErr
}

func getRequestID(ctx *gin.Context) string {
	if reqID, exists := ctx.Get("requestId"); exists {
		if id, ok := reqID.(string); ok {
			return id
		}
	}
	return ""
}

func getDocumentationURL(ctx *gin.Context) string {
	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}

	host := ctx.Request.Host
	fullPath := strings.TrimPrefix(ctx.FullPath(), "/api/")
	pathSegments := strings.Split(fullPath, "/")

	if len(pathSegments) == 1 && pathSegments[0] != "" {
		fullPath = fmt.Sprintf("%s/%s", fullPath, strings.ToLower(ctx.Request.Method))
	}

	return fmt.Sprintf("%s://%s/api/docs/#/%s", scheme, host, fullPath)
}

func getCurrentTimestamp() string {
	loc := time.FixedZone("WIB", 7*60*60)
	return time.Now().In(loc).Format("2006-01-02 15:04:05.000")
}

func FailedResponse(ctx *gin.Context, statusCode int, fallback model.Message, err error) model.WebResponse[any] {
	if resolvedCode := GetHTTPStatusCode(err); resolvedCode != http.StatusOK {
		statusCode = resolvedCode
	}

	var msg model.Message
	switch {
	case err == nil:
		msg = fallback
	case isMultilangError(err.Error()):
		msg = ParseMultilangError(err)
	default:
		shortErr := extractErrorMessage(err.Error())
		msg = model.Message{"en": shortErr, "id": shortErr}
	}

	return model.WebResponse[any]{
		Status:           "error",
		StatusCode:       statusCode,
		Message:          msg,
		RequestID:        getRequestID(ctx),
		Timestamp:        getCurrentTimestamp(),
		Path:             ctx.FullPath(),
		DocumentationURL: getDocumentationURL(ctx),
	}
}

func SuccessResponse[T any](ctx *gin.Context, statusCode int, message model.Message, data T) model.WebResponse[T] {
	return model.WebResponse[T]{
		Status:           "success",
		StatusCode:       statusCode,
		Message:          message,
		Data:             data,
		RequestID:        getRequestID(ctx),
		Timestamp:        getCurrentTimestamp(),
		Path:             ctx.Request.URL.Path,
		DocumentationURL: getDocumentationURL(ctx),
	}
}

func SuccessWithPaginationResponse[T any](
	ctx *gin.Context,
	statusCode int,
	message model.Message,
	data []T,
	paging model.PageMetadata,
	documentationURL ...string,
) model.WebResponse[[]T] {
	return model.WebResponse[[]T]{
		Status:           model.StatusSuccess,
		StatusCode:       statusCode,
		Message:          message,
		Data:             data,
		Paging:           &paging,
		RequestID:        getRequestID(ctx),
		Timestamp:        getCurrentTimestamp(),
		Path:             ctx.Request.URL.Path,
		DocumentationURL: getDocumentationURL(ctx),
	}
}
