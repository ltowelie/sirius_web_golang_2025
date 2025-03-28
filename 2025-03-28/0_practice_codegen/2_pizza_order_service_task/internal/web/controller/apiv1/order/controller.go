package order

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"pizza_order_service/internal/web/api"
)

type Controller struct {
}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) DefineRoutes(r gin.IRouter) {
	api.RegisterHandlers(r, c)
}

func (c *Controller) GetOrders(ctx *gin.Context, params api.GetOrdersParams) {
	// Второй способ валидации - не через middleware
	// но нужно использовать раcширение и руками дополнять спецификацию, что не очень удобно
	//
	// При каждом запросе инициализировать валидадор - плохое решение
	// тут показано, просто для примера. Как думаете, где лучше его инициализировать?
	//validate := validator.New()
	//err := validate.Struct(params)
	//if err != nil {
	//	errValidation := err.(validator.ValidationErrors)
	//
	//	ctx.JSON(http.StatusBadRequest, api.ErrorResponse{})
	//}

	// В params приходят параметры из запроса (но их также можно получить из контекста, но
	// все же лучше из сгенерированного кода, так как там все параметры уже преобразованы)
	lp := make([]any, 0)
	if params.Status != nil {
		lp = append(lp, slog.Any("status", *params.Status))
	}
	if params.Limit != nil {
		lp = append(lp, slog.Any("limit", *params.Limit))
	}
	if params.Offset != nil {
		lp = append(lp, slog.Any("offset", *params.Offset))
	}
	if params.Sort != nil {
		lp = append(lp, slog.Any("sort", *params.Sort))
	}
	slog.Debug("request params",
		lp...,
	)

	// Здесь пишем нашу бизнес логику
	// Например вызов метода сервиса, который в свою очередь вызовет метод репозитория
	// для получения данных из хранилища.
	// Имитация данных
	mockUUID := uuid.New()

	// Обратите внимание - используем сгенерированные типы из пакета api
	orders := &[]api.Order{
		{Id: mockUUID, Status: api.OrderStatusPending},
	}

	// Отправляем ответ через Gin
	// Здесь тоже используем сгенерированные типы из пакета api
	ctx.JSON(http.StatusOK, api.OrdersResponse{
		Data: orders,
	})
}
