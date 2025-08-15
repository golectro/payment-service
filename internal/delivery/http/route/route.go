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
	api := c.App.Group("/api")

	c.RegisterCommonRoutes(c.App)
	c.RegisterSwaggerRoutes(api)
	c.RegisterPaymentRoutes(api)
}
