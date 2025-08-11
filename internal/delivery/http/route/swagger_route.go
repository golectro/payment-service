package route

import (
	"github.com/gin-gonic/gin"
)

func (c *RouteConfig) RegisterSwaggerRoutes(rg *gin.RouterGroup) {
	swagger := rg.Group("/docs")
	swagger.StaticFile("/swagger.json", "./docs/swagger.json")
	swagger.GET("/", c.SwaggerController.SwaggerDocHandler)
}
