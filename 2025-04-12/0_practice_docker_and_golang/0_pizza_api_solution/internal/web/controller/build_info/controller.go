package build_info

import (
	"github.com/gin-gonic/gin"

	"pizza_api/internal/app/build"
)

type Controller struct {
	info *build.Info
}

func NewController(bi *build.Info) *Controller {
	return &Controller{info: bi}
}

func (c *Controller) DefineRoutes(r gin.IRouter) {
	r.GET("/build_info", c.BuildInfo)
}

func (c *Controller) BuildInfo(ctx *gin.Context) {
	ctx.JSON(200, c.info)
}
