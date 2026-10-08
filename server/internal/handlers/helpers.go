package handlers

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidJSON = errors.New("invalid request body")
	ErrInvalidParam = errors.New("invalid path parameter")
	ErrMissingUserID = errors.New("user id not in context")
)


func getUserID(c *gin.Context) (int64, error) {
	value, exists := c.Get("userID")
	if !exists {
		return 0, ErrMissingUserID
	}

	userID, ok := value.(int64)
	if !ok {
		return 0, ErrMissingUserID
	}

	return userID, nil
}