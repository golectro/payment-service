package route

import (
	"github.com/gin-gonic/gin"
)

func (c *RouteConfig) RegisterSwaggerRoutes(app *gin.Engine) {
	swagger := app.Group("/docs")
	swagger.StaticFile("/swagger.json", "./docs/swagger.json")
	swagger.GET("/", c.SwaggerController.SwaggerDocHandler)
}
