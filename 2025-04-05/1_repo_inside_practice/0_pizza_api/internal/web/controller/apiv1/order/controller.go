package order

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"repository_example/internal/models"
	"repository_example/internal/web/controller/apiv1"
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
