package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/qullDev/book_API/internal/config"
	appauth "github.com/qullDev/book_API/internal/pkg/auth"
)

// NewJWTAuth memvalidasi Authorization Bearer token menggunakan helper JWT
func NewJWTAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "missing or invalid Authorization header"})
			return
		}

		claims, err := appauth.ParseToken(cfg, parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid or expired token"})
			return
		}

		// simpan userID di context untuk digunakan handler
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
