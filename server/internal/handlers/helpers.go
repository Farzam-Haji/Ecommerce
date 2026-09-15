package handlers

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidJSON = errors.New("Couldnt Bind Request")
	ErrMissingUserID = errors.New("User ID Not In Context")
	ErrInvalidParam = errors.New("Param is Invalid")
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