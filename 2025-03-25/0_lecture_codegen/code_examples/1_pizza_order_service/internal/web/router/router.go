package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/oapi-codegen/gin-middleware"

	"pizza_order_service/internal/web/middleware/oapi_validator"
)

func New() (*gin.Engine, error) {
	r := gin.New()

	oapiV, err := oapi_validator.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize gin middleware - OapiValidator: %w", err)
	}

	r.Use(gin.Recovery(), gin.Logger(), oapiV)

	return r, nil
}
