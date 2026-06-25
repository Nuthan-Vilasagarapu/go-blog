package routers

import (
	"github.com/gin-gonic/gin"
)

func LoadRoutes(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")

	LoadAuthRoutes(router)
	LoadBlogRoutes(router)
	LoadViewRoutes(router)
}
