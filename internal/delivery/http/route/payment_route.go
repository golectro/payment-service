package route

import (
	"github.com/gin-gonic/gin"
)

func (c *RouteConfig) RegisterPaymentRoutes(rg *gin.RouterGroup) {
	payment := rg.Group("payment")

	payment.POST("/invoice", c.AuthMiddleware, c.PaymentController.CreateInvoice)
	payment.GET("/invoice", c.AuthMiddleware, c.PaymentController.GetInvoice)
	payment.POST("/xendit/callback", c.PaymentController.XenditCallback)
}
