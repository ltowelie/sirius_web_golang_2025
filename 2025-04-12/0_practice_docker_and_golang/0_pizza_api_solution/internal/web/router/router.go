package router

import (
	"github.com/gin-gonic/gin"
)

func New() (*gin.Engine, error) {
	r := gin.New()

	r.Use(gin.Logger(), gin.Recovery())

	return r, nil
}
