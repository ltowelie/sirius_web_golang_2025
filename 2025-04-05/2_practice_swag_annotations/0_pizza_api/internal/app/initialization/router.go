package initialization

import (
	"github.com/gin-gonic/gin"

	"pizza_api/internal/service"
	"pizza_api/internal/web/controller/apiv1/order"
	order2 "pizza_api/internal/web/controller/swagger"
	"pizza_api/internal/web/router"
)

type Controller interface {
	DefineRoutes(r gin.IRouter)
}

func InitializeRouter(s *service.Services) (*gin.Engine, error) {
	r, err := router.New()
	if err != nil {
		return nil, err
	}

	cs := InitControllers(s)
	for _, c := range cs {
		c.DefineRoutes(r)
	}

	return r, nil
}

func InitControllers(s *service.Services) []Controller {
	return []Controller{
		order.NewController(s.Orders),
		order2.NewController(),
	}
}
