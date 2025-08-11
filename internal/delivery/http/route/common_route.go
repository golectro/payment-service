package route

import (
	"golectro-payment/internal/constants"
	"golectro-payment/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (c *RouteConfig) RegisterCommonRoutes(app *gin.Engine) {
	welcomeHandler := func(ctx *gin.Context) {
		res := utils.SuccessResponse(ctx, http.StatusOK, constants.WelcomeMessage, "Welcome to Golectro Payment API")
		ctx.JSON(res.StatusCode, res)
	}

	app.GET("/", welcomeHandler)
	app.GET("/api", welcomeHandler)

	app.NoRoute(func(ctx *gin.Context) {
		res := utils.FailedResponse(ctx, http.StatusNotFound, constants.NotFound, nil)
		ctx.AbortWithStatusJSON(res.StatusCode, res)
	})
}
