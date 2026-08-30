package handlers

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/btalv/Ecommerce/server/internal/models"
	"github.com/btalv/Ecommerce/server/internal/services"
	"github.com/gin-gonic/gin"
)

type CartServiceAPI interface {
	GetUserCartItems(ctx context.Context, userID int64) ([]*models.CartItem, []string, error)
	AddToCart(ctx context.Context, userID int64, productID int64, pointerQuantity *int32) (error)
	UpdateItemQuantity(ctx context.Context, userID int64, itemID int64, quantity int32) (error)
	DeleteFromCart(ctx context.Context, userID int64, itemID int64) (error)
}

type CartHandler struct {
	Service  	CartServiceAPI
}

func NewCartHandler(serviceInterface CartServiceAPI) *CartHandler {
	return &CartHandler{
		Service: serviceInterface,
	}
}

func (ch *CartHandler) GetMyCartItems(c *gin.Context) {
	value, exists := c.Get("userID")
	if !exists {
		c.Error(services.ErrInternal)
		return
	}
	userID, ok := value.(int64)
	if !ok {
		c.Error(services.ErrInternal)
		return
	}

	items, deleted, err := ch.Service.GetUserCartItems(c.Request.Context(), userID)
	if errors.Is(err, services.ErrCartNotFound) {
		c.JSON(200, "your cart is empty")
		return
	}
	if err != nil {
		c.Error(fmt.Errorf("%w: %v", services.ErrInternal, err))
		return
	}

	response := gin.H{
	    "items": items,
	}

	if len(deleted) != 0 {
	    response["notification"] = "some items had been removed from your cart"
	    response["deleted_items"] = deleted
	}

	c.JSON(200, response)
}

func (ch *CartHandler) AddItemToMyCart(c *gin.Context) {
	var request AddItemRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(fmt.Errorf("%w: %v", ErrInvalidJSON, err))
		return
	}

	value, exists := c.Get("userID")
	if !exists {
		c.Error(services.ErrInternal)
		return
	}
	userID, ok := value.(int64)
	if !ok {
		c.Error(services.ErrInternal)
		return
	}

	err := ch.Service.AddToCart(c.Request.Context(), userID, request.ProductID, request.Quantity)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(201, gin.H{"status": "created"})
}

func (ch *CartHandler) UpdateMyItemQuantity(c *gin.Context) {
	param := c.Param("itemID")
	itemID, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		c.Error(services.ErrInvalidID)
		return
	}

	var request UpdateItemQuantityRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(fmt.Errorf("%w: %v", ErrInvalidJSON, err))
		return
	}

	value, exists := c.Get("userID")
	if !exists {
		c.Error(services.ErrInternal)
		return
	}
	userID, ok := value.(int64)
	if !ok {
		c.Error(services.ErrInternal)
		return
	}

	err = ch.Service.UpdateItemQuantity(c.Request.Context(), userID, itemID, request.Quantity)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, gin.H{"status": "success"})
}

func (ch *CartHandler) DeleteMyItem(c *gin.Context) {
	param := c.Param("itemID")
	itemID, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		c.Error(services.ErrInvalidID)
		return
	}

	value, exists := c.Get("userID")
	if !exists {
		c.Error(services.ErrInternal)
		return
	}
	userID, ok := value.(int64)
	if !ok {
		c.Error(services.ErrInternal)
		return
	}

	err = ch.Service.DeleteFromCart(c.Request.Context(), userID, itemID)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, gin.H{"status": "success"})
}
