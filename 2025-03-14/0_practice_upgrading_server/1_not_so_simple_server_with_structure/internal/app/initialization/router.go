package initialization

import (
	"net/http"

	"practice/internal/web/controller/apiv1/home_work"
	"practice/internal/web/controller/apiv1/user"
	"practice/internal/web/router"
)

type Controller interface {
	DefineRoutes(mux *http.ServeMux)
}

func NewRouter() *http.ServeMux {
	r := router.New()

	cs := InitControllers()
	for _, c := range cs {
		c.DefineRoutes(r)
	}

	return r
}

func InitControllers() []Controller {
	return []Controller{
		home_work.NewController(),
		user.NewController(),
	}
}
