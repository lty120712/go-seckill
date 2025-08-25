package middleware

import (
	"github.com/gin-gonic/gin"
	"go-chat/internal/utils/jwtUtil"
	"net/http"
	"strings"
)

// AuthMiddleware 用于验证 JWT Token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")

		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(tokenStr, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}

		tokenStr = tokenStr[len("Bearer "):]

		claims, err := jwtUtil.ParseJWT(tokenStr)
		if err != nil {

			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}
		c.Set("id", claims.ID)
		c.Next()
	}
}
