package swagger

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"pizza_api/api_docs"
)

type Controller struct {
}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) DefineRoutes(r gin.IRouter) {
	_ = api_docs.SwaggerInfo

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
