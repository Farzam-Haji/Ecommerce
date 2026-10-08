package middlewares

import (
	"errors"
	"strings"

	"github.com/btalv/Ecommerce/server/internal/auth"
	"github.com/btalv/Ecommerce/server/internal/services"
	"github.com/gin-gonic/gin"
)

var ErrUnauthorized = errors.New("authentication required")

func Authentication(jwtManager	*auth.JWTManager) gin.HandlerFunc {
	return	func(c *gin.Context) {
		authJWT := c.Request.Header.Get("Authorization")
		if authJWT == "" {
			c.Next()
			return 
		}
		if !strings.HasPrefix(authJWT, "Bearer ") {
			c.Error(ErrUnauthorized)
			c.Abort()
			return 
		}
		authToken := strings.TrimPrefix(authJWT, "Bearer ")
		claims, err := jwtManager.ValidateToken(authToken)
		if err != nil {
			c.Error(err)
			c.Abort()
			return 
		}
		c.Set("claims", claims)

		c.Next()		
	}
}

func RequireUser(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context){
		value, exists := c.Get("claims")
		if !exists {
			c.Error(ErrUnauthorized)
			c.Abort()
			return 
		}
		claims, ok := value.(*auth.CustomClaims)
		if !ok {
			c.Error(services.ErrInternal)
			c.Abort()
			return 
		}
		c.Set("userID", claims.UserID)

		c.Next()
	}
}