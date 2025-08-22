package route

import (
	"golectro-payment/internal/delivery/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type RouteConfig struct {
	App               *gin.Engine
	AuthMiddleware    gin.HandlerFunc
	Viper             *viper.Viper
	SwaggerController *http.SwaggerController
	PaymentController *http.PaymentController
}

func (c *RouteConfig) Setup() {
	api := c.App.Group("/api/v1")

	c.RegisterPaymentRoutes(api)
	c.RegisterSwaggerRoutes(c.App)
	c.RegisterCommonRoutes(c.App)
}
