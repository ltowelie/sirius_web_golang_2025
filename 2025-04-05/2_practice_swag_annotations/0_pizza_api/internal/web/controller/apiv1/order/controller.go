package order

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"pizza_api/internal/models"
	"pizza_api/internal/web/controller/apiv1"
)

type OrdersService interface {
	GetByID(ctx context.Context, id int) (*models.Order, error)
	Create(ctx context.Context, o *models.Order) (*models.Order, error)
}

type Controller struct {
	orders OrdersService
}

func NewController(o OrdersService) *Controller {
	return &Controller{orders: o}
}

func (c *Controller) DefineRoutes(r gin.IRouter) {
	g := r.Group(apiv1.Group)

	g.GET("/orders/:id", c.GetOrderByID)
	g.POST("/orders", c.PostOrder)
}

// GetOrderByID godoc
//
//	@Summary		Получить заказ по ID
//	@Description	Получает информацию о заказе по его идентификатору
//	@Tags			orders
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"ID заказа"
//	@Success		200	{object}	models.Order
//	@Failure		400	{object}	object{error=string}	"Неверный формат ID"
//	@Failure		500	{object}	object{error=string}	"Ошибка сервера"
//	@Router			/v1/orders/{id} [get]
func (c *Controller) GetOrderByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		errStr := "failed to parse id: " + err.Error()
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errStr})

		return
	}

	o, err := c.orders.GetByID(ctx, int(id))
	if err != nil {
		errStr := "failed to get order by id: " + err.Error()
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": errStr})

		return
	}

	ctx.JSON(http.StatusOK, o)
}

// PostOrder godoc
//
//	@Summary		Создать новый заказ
//	@Description	Создает новый заказ
//	@Tags			orders
//	@Accept			json
//	@Produce		json
//	@Param			order	body		models.Order	true	"Заказ"
//	@Success		201		{object}	models.Order
//	@Failure		400		{object}	object{error=string}	"Неверный формат данных"
//	@Failure		500		{object}	object{error=string}	"Ошибка сервера"
//	@Router			/v1/orders [post]
func (c *Controller) PostOrder(ctx *gin.Context) {
	order := &models.Order{}
	err := ctx.BindJSON(order)
	if err != nil {
		errStr := "failed to parse order from JSON: " + err.Error()
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errStr})

		return
	}

	order, err = c.orders.Create(ctx.Request.Context(), order)
	if err != nil {
		errStr := "failed to create order: " + err.Error()
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": errStr})
	}

	ctx.JSON(http.StatusCreated, order)
}
