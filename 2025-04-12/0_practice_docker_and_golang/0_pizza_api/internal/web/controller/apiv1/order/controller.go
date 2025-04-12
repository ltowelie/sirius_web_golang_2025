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
	Update(ctx context.Context, order *models.Order) (*models.Order, error)
	Delete(ctx context.Context, id int) error
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
	g.PUT("/orders/:id", c.UpdateOrder)
	g.DELETE("/orders/:id", c.DeleteOrder)
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

// UpdateOrder godoc
//
//	@Summary		Обновить заказ
//	@Description	Обновляет существующий заказ по ID
//	@Tags			orders
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int				true	"ID заказа"
//	@Param			order	body		models.Order	true	"Обновленные данные"
//	@Success		204
//	@Failure		400		{object}	object{error=string}	"Неверный запрос"
//	@Failure		404		{object}	object{error=string}	"Заказ не найден"
//	@Failure		500		{object}	object{error=string}	"Ошибка сервера"
//	@Router			/api/v1/orders/{id} [put]
func (c *Controller) UpdateOrder(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	var order models.Order
	err = ctx.BindJSON(&order)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if order.ID != 0 && order.ID != id {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "id mismatch"})
		return
	}
	order.ID = id

	orderU, err := c.orders.Update(ctx.Request.Context(), &order)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": "failed to update order: " + err.Error()},
		)

		return
	}

	ctx.JSON(http.StatusOK, orderU)
}

// DeleteOrder godoc
//
//	@Summary		Удалить заказ
//	@Description	Удаляет заказ по ID
//	@Tags			orders
//	@Produce		json
//	@Param			id	path		int	true	"ID заказа"
//	@Success		204
//	@Failure		400	{object}	object{error=string}	"Неверный ID"
//	@Failure		404	{object}	object{error=string}	"Заказ не найден"
//	@Failure		500	{object}	object{error=string}	"Ошибка сервера"
//	@Router			/api/v1/orders/{id} [delete]
func (c *Controller) DeleteOrder(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})

		return
	}

	if err := c.orders.Delete(ctx.Request.Context(), id); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"error": "failed to delete order: " + err.Error()},
		)

		return
	}

	ctx.Status(http.StatusNoContent)
}
