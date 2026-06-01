package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
)


func HeaderCheck() gin.HandlerFunc {
	return	func(c *gin.Context) {
		contentType := c.GetHeader("Content-Type")
		if !strings.HasPrefix(contentType, "application/json") {
			c.AbortWithStatus(415)
			return 
		}

		c.Next()
	}
}