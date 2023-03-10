package middleware

import (
	"net/http"
	"oprec/go-blog/service"
	"oprec/go-blog/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response := utils.BuildErrorResponse("No token found", http.StatusUnauthorized)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		if !strings.Contains(authHeader, "Bearer ") {
			response := utils.BuildErrorResponse("No token found", http.StatusUnauthorized)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		authHeader = strings.Replace(authHeader, "Bearer ", "", -1)

		tokenService := service.NewJWTService()

		token, err := tokenService.ValidateToken(authHeader)
		if err != nil {
			response := utils.BuildErrorResponse("Invalid token", http.StatusUnauthorized)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		if !token.Valid {
			response := utils.BuildErrorResponse("Invalid token", http.StatusUnauthorized)
			c.AbortWithStatusJSON(http.StatusForbidden, response)
			return
		}

		c.Next()
	}
}
