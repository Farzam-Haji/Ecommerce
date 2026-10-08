package handlers

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/btalv/Ecommerce/server/internal/models"
	"github.com/gin-gonic/gin"
)

type OrderServiceAPI interface {
	CreateOrder(ctx context.Context, userID int64) (*models.Order, error)
	GetUserOrders(ctx context.Context, userID int64) ([]*models.Order, error)
	GetUserOrderByID(ctx context.Context, orderID int64, userID int64) (*models.Order, error)
	GetOrderItems(ctx context.Context, orderID int64) ([]*models.OrderItem, error)
	CreatePayment(ctx context.Context, orderID int64) (string, error)
	ProcessPayment(ctx context.Context, orderID int64, paymentID int64, status string) (error)
	CancelUserOrder(ctx context.Context, orderID int64, userID int64) (error)
}

type OrderHandler struct {
	service 	OrderServiceAPI
	serverURL 	string
}

func NewOrderHandler(serviceInterface OrderServiceAPI, serverURL string) (*OrderHandler) {
	return &OrderHandler{
		service: serviceInterface,
		serverURL: serverURL,
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

	paymentURL, err := oh.service.CreatePayment(c.Request.Context(), order.ID)
	if err != nil {
		c.Error(err)
		return
	}

	// c.Redirect(303, paymentURL)

	c.JSON(200, gin.H{
		"order": order,
		"payment_url": paymentURL,
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

func (oh *OrderHandler) Result(c *gin.Context) {
	var request PaymentResultRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(fmt.Errorf("%w: %w", ErrInvalidJSON, err))
		return
	}

	log.Printf(
        "payment result: payment=%d order=%d status=%s",
        request.PaymentID,
        request.OrderID,
        request.Status,
    )

	err := oh.service.ProcessPayment(c.Request.Context(), request.OrderID, request.PaymentID, request.Status)
	if err != nil {
		c.Error(err)
		return
	}



	c.Status(200)
}

func (oh *OrderHandler) CancelOrder(c *gin.Context) {
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

	err = oh.service.CancelUserOrder(c.Request.Context(), orderID, userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(200)
}