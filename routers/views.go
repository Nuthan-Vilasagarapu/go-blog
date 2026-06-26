package routers

import (
	"go-blog/controllers"

	"github.com/gin-gonic/gin"
)

func LoadViewRoutes(router *gin.Engine) {
	router.GET("/", controllers.HomePage)
	router.Any("/health", controllers.Health)
}
