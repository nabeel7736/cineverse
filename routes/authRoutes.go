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

	rg.GET("/movies", middlewares.AuthMiddleware(), controllers.GetAllMovies(db))

}
