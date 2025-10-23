package routes

import (
	"cineverse/controllers"
	"cineverse/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AdminRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	rg.GET("/admin/dashboard", middlewares.AdminMiddleware(db), controllers.AdminDashboard(db))
}
