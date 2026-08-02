package middlewares

import (
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
	"time"

	"github.com/btalv/Ecommerce/server/internal/auth"
	"github.com/btalv/Ecommerce/server/internal/repositories"
	"github.com/btalv/Ecommerce/server/internal/services"
	"github.com/gin-gonic/gin"
)


func CustomLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		c.Next()

		request, ok := c.MustGet("requestID").(string)
		if !ok {
			c.AbortWithError(500, errors.New("something went wrong with request ID"))
			return 
		}
		status := c.Writer.Status()
		duration := time.Since(start)
		log.Printf("request '%s' returned '%v' and took '%s'", request, status, duration)
	}
}

func RequestID() gin.HandlerFunc {
	return func (c *gin.Context) {
		reqID := 10000 + rand.IntN(89999)
		reqString := fmt.Sprintf("req-%v", reqID)
		c.Set("requestID", reqString)

		c.Next()
	}
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Writer.Written() {
		    return
		}

		if len(c.Errors) == 0 {
			return 
		}
		err := c.Errors.Last().Err
		switch{
		case errors.Is(err, repositories.ErrDuplicateKey):
			c.JSON(409, gin.H{"error": repositories.ErrDuplicateKey.Error()})
		case errors.Is(err, services.ErrInvalidID):
			c.JSON(400, gin.H{"error": services.ErrInvalidID.Error()})
		case errors.Is(err, services.ErrInvalidName):
			c.JSON(400, gin.H{"error": services.ErrInvalidName.Error()})
		case errors.Is(err, services.ErrInvalidQuantity):
			c.JSON(400, gin.H{"error": services.ErrInvalidQuantity.Error()})
		case errors.Is(err, services.ErrProductNotFound):
			c.JSON(404, gin.H{"error": services.ErrProductNotFound.Error()})
		case errors.Is(err, auth.ErrInvalidToken):
			log.Println(err)
			c.JSON(401, gin.H{"error": auth.ErrInvalidToken.Error()})
		case errors.Is(err, auth.ErrExpiredToken):
			c.JSON(401, gin.H{"error": auth.ErrExpiredToken.Error()})
		case errors.Is(err, ErrUnauthorized):
			c.JSON(401, gin.H{"error": ErrUnauthorized.Error()})
		default:
			log.Println(err)
			c.JSON(500, gin.H{"error": services.ErrInternal.Error()})
		}
	}
}