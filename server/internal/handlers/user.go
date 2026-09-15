package handlers

import (
	"context"
	"fmt"

	"github.com/btalv/Ecommerce/server/internal/models"
	"github.com/gin-gonic/gin"
)

type UserServiceAPI interface {
	Register(context.Context, string, string)(*models.User, string, error)
	Login(context.Context, string, string)(string, error)
}

type UserHandler struct {
	Service 	UserServiceAPI
}

func NewUserHandler(serviceInterface UserServiceAPI) *UserHandler {
	return &UserHandler{
		Service: serviceInterface,
	}
}

func (uh *UserHandler) Register(c *gin.Context) {
	var request RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(fmt.Errorf("%w: %w", ErrInvalidJSON, err))
		return
	}

	user, jwt, err := uh.Service.Register(c.Request.Context(), request.Email, request.Password)
	if err != nil {
		c.Error(err)
		return
	}
	
	c.JSON(201, gin.H{
		"status": "created",
		"user": user.Email,
		"token": jwt,
	})
}

func (uh *UserHandler) Login(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(fmt.Errorf("%w: %w", ErrInvalidJSON, err))
		return
	}

	jwt, err := uh.Service.Login(c.Request.Context(), request.Email, request.Password)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"token": jwt,
	})
}