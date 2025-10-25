package middlewares

import (
	"cineverse/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			// Try reading from cookie if header missing
			cookie, err := ctx.Cookie("access_token")
			if err != nil || cookie == "" {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token missing"})
				ctx.Abort()
				return
			}
			token = cookie
		}

		if strings.HasPrefix(token, "Bearer ") {
			token = strings.TrimPrefix(token, "Bearer ")
		}

		claims, err := utils.ValidateAccessToken(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			ctx.Abort()
			return
		}

		// userID := uint(claims["user_id"].(float64))

		// ctx.Set("user_id", userID)
		ctx.Set("claims", claims)
		ctx.Next()
	}
}
