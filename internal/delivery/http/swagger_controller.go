package http

import (
	"golectro-payment/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type SwaggerController struct {
	Log *logrus.Logger
}

func NewSwaggerController(log *logrus.Logger) *SwaggerController {
	return &SwaggerController{
		Log: log,
	}
}

func (c *SwaggerController) SwaggerDocHandler(ctx *gin.Context) {
	ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(utils.SwaggerHTML()))
}
