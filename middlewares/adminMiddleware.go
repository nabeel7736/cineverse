// package middlewares

// import (
// 	"cineverse/models"
// 	"cineverse/utils"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"gorm.io/gorm"
// )

// func AdminMiddleware(db *gorm.DB) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		token, err := c.Cookie("access_token")
// 		if err != nil || token == "" {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token missing"})
// 			c.Abort()
// 			return
// 		}

// 		claims, err := utils.ValidateAccessToken(token)
// 		if err != nil {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
// 			c.Abort()
// 			return
// 		}

// 		userID := uint(claims["user_id"].(float64))
// 		var user models.User
// 		if err := db.First(&user, userID).Error; err != nil {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
// 			c.Abort()
// 			return
// 		}

// 		if user.Role != "admin" {
// 			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
// 			c.Abort()
// 			return
// 		}

// 		// c.Set("user_id", userID)
// 		c.Next()
// 	}
// }

package middlewares

import (
	"cineverse/models"
	"cineverse/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AdminMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from cookie
		token, err := c.Cookie("access_token")
		if err != nil || token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token missing"})
			c.Abort()
			return
		}

		// Validate token
		claims, err := utils.ValidateAccessToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Ensure claims is not nil
		if claims == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Access UserID from claims struct (fix for line 28)
		userID := claims.UserID
		if userID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: missing user_id"})
			c.Abort()
			return
		}

		// Fetch user from database
		var user models.User
		if err := db.First(&user, userID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			c.Abort()
			return
		}

		// Check if user has admin role
		if user.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: admin role required"})
			c.Abort()
			return
		}

		// Set userID or claims in context for downstream handlers
		c.Set("user_id", userID)
		c.Set("claims", claims) // Optional: store full claims for flexibility

		// Proceed to next handler
		c.Next()
	}
}
