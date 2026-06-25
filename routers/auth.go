package routers

import (
	"go-blog/controllers"
	"go-blog/middlewares"

	"github.com/gin-gonic/gin"
)

func LoadAuthRoutes(router *gin.Engine) {
	// API
	router.POST("/auth/register", controllers.RegisterUser)
	router.POST("/auth/login", controllers.LoginUser)
	router.GET("/users", middlewares.AuthMiddleware(), controllers.GetUsers)

	// HTML
	router.GET("/login", controllers.LoginPage)
}
