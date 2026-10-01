package middlewares

import (
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
	"time"

	"github.com/btalv/Ecommerce/server/internal/auth"
	"github.com/btalv/Ecommerce/server/internal/bank"
	"github.com/btalv/Ecommerce/server/internal/handlers"
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
		//auth
		case errors.Is(err, auth.ErrInvalidToken):
			log.Println(err)
			c.JSON(401, gin.H{"error": auth.ErrInvalidToken.Error()})
		case errors.Is(err, auth.ErrExpiredToken):
			c.JSON(401, gin.H{"error": err.Error()})
		case errors.Is(err, ErrUnauthorized):
			c.JSON(401, gin.H{"error": err.Error()})

		//handler
		case errors.Is(err, handlers.ErrInvalidJSON):
			log.Println(err)
			c.JSON(400, gin.H{"error": handlers.ErrInvalidJSON.Error()})
		case errors.Is(err, handlers.ErrInvalidParam):
			log.Println(err)
			c.JSON(400, gin.H{"error": handlers.ErrInvalidParam.Error()})
		case errors.Is(err, handlers.ErrMissingUserID):
			log.Println(err)
			c.JSON(401, gin.H{"error": handlers.ErrMissingUserID.Error()})

		//user
		case errors.Is(err, services.ErrEmailUnavailable):
			c.JSON(400, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrInvalidPassword):
			c.JSON(400, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrWrongCredential):
			c.JSON(401, gin.H{"error": err.Error()})
		case errors.Is(err,services.ErrTokenGenFailed):
			log.Println(err)
			c.JSON(500, gin.H{"error": services.ErrTokenGenFailed.Error()})

		//service general
		case errors.Is(err, services.ErrInvalidID):
			c.JSON(400, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrInvalidName):
			c.JSON(400, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrInvalidQuantity):
			c.JSON(400, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrInvalidPrice):
			c.JSON(400, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrUnauthorized):
			c.JSON(401, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrInternal):
			log.Println(err)
			c.JSON(500, gin.H{"error": services.ErrInternal.Error()})
		
		//service product
		case errors.Is(err, services.ErrProductNotFound):
			c.JSON(404, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrDuplicateProduct):
			c.JSON(409, gin.H{"error": err.Error()})


		//service cart and order
		case errors.Is(err, services.ErrInsufficientQuantity):
			c.JSON(409, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrCartNotFound):
			c.JSON(404, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrItemNotFound):
			c.JSON(404, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrEmptyOrder):
			c.JSON(400, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrOrderNotFound):
			c.JSON(404, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrNonPending):
			c.JSON(404, gin.H{"error": err.Error()})
		
		//bank
		case errors.Is(err, bank.ErrRequestBank):
			log.Println(err)
			c.JSON(500, gin.H{"error": bank.ErrRequestBank.Error()})

		default:
			log.Println(err)
			c.JSON(500, gin.H{"error": services.ErrInternal.Error()})
		}
	}
}