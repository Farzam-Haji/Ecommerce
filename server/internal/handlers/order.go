package handlers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/btalv/Ecommerce/server/internal/models"
	"github.com/gin-gonic/gin"
)

type OrderServiceAPI interface {
	CreateOrder(ctx context.Context, userID int64) (*models.Order, error)
	GetUserOrders(ctx context.Context, userID int64) ([]*models.Order, error)
	GetUserOrderByID(ctx context.Context, orderID int64, userID int64) (*models.Order, error)
	GetOrderItems(ctx context.Context, orderID int64) ([]*models.OrderItem, error)
}

type OrderHandler struct {
	service 	OrderServiceAPI
}

func NewOrderHandler(serviceInterface OrderServiceAPI) (*OrderHandler) {
	return &OrderHandler{
		service: serviceInterface,
	}
}

func (oh *OrderHandler) CreateOrder(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	order, err := oh.service.CreateOrder(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(201, gin.H{
		"status": "success",
		"order": order,
	})
}

func (oh *OrderHandler) GetMyOrders(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	orders, err := oh.service.GetUserOrders(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{
		"orders": orders,
	})
}

func (oh *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	orderID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.Error(fmt.Errorf("%w: %w", ErrInvalidParam, err))
		return
	}

	userID, err := getUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	order, err := oh.service.GetUserOrderByID(c.Request.Context(), orderID, userID)
	if err != nil {
		c.Error(err)
		return
	}

	orderItems, err := oh.service.GetOrderItems(c.Request.Context(), order.ID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{
		"order": order,
		"items": orderItems,
	})
}