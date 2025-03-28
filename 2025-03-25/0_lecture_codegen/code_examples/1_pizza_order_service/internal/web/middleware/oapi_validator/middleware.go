package oapi_validator

import (
	"github.com/gin-gonic/gin"
	ginmiddleware "github.com/oapi-codegen/gin-middleware"

	"pizza_order_service/internal/web/api"
)

func New() (gin.HandlerFunc, error) {
	sw, err := api.GetSwagger()
	if err != nil {
		return nil, err
	}
	m := ginmiddleware.OapiRequestValidator(sw)

	return m, nil
}
