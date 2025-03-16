package user

import "net/http"

type Controller struct {
}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) DefineRoutes(mux *http.ServeMux) {
}
