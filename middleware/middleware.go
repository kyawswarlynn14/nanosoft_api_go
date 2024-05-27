package middleware

import (
	"net/http"

	token "nanosoft/tokens"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		ClientToken := c.Request.Header.Get("token")
		if ClientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No Authorization Header Provided"})
			c.Abort()
			return
		}
		claims, err := token.ValidateToken(ClientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Set("uid", claims.Uid)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func AuthorizeRole(roles []int) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Role not found in context"})
			c.Abort()
			return
		}

		userRoleInt, ok := userRole.(int)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Role is not an integer"})
			c.Abort()
			return
		}

		roleValid := false
		for _, role := range roles {
			if userRoleInt == role {
				roleValid = true
				break
			}
		}

		if !roleValid {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have the necessary permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}
