package initialization

import (
	"github.com/gin-gonic/gin"

	"pizza_order_service/internal/web/controller/apiv1/order"
	"pizza_order_service/internal/web/router"
)

type Controller interface {
	DefineRoutes(r gin.IRouter)
}

func InitializeRouter() (*gin.Engine, error) {
	r, err := router.New()
	if err != nil {
		return nil, err
	}

	cs := InitControllers()
	for _, c := range cs {
		c.DefineRoutes(r)
	}

	return r, nil
}

func InitControllers() []Controller {
	return []Controller{
		order.NewController(),
	}
}
