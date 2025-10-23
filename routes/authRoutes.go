package routes

import (
	"cineverse/controllers"
	"cineverse/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthRoute(rg *gin.RouterGroup, db *gorm.DB) {
	rg.GET("/register", controllers.ShowRegisterpage)
	rg.POST("/register", controllers.RegisterHandler(db))

	rg.GET("/login", controllers.ShowLoginPage)
	rg.POST("/login", controllers.LoginHandler(db))
	rg.POST("/refresh", controllers.RefreshTokenHandler(db))

	rg.POST("/logout", middlewares.AuthMiddleware(), controllers.LogoutHandler(db))

	rg.GET("/movies", middlewares.AuthMiddleware(), func(ctx *gin.Context) {

		userID, _ := ctx.Get("user_id")
		ctx.HTML(200, "movies/book.html", gin.H{
			"title":      "Profile",
			"user_id":    userID,
			"ActivePage": "profile",
		})
	})

}
