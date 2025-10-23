package routes

import (
	// "cineverse/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// load templates
	r.LoadHTMLGlob("templates/**/*.html")

	// static files (if any)
	r.Static("/static", "./static")

	api := r.Group("/")
	AuthRoute(api, db)

	// other route groups will be added here (movies, admin, booking etc.)
	return r
}
