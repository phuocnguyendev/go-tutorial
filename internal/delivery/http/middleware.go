package httpdelivery

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"go-tutorial/internal/usecase"
)

func AuthMiddleware(uc usecase.AuthUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid Authorization header"})
			return
		}

		tokenStr := strings.TrimPrefix(h, "Bearer ")

		userID, err := uc.ParseToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
