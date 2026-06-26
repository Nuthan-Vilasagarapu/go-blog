package routers

import (
	"github.com/gin-gonic/gin"
)

func LoadRoutes(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	LoadViewRoutes(router)
	LoadAuthRoutes(router)
	LoadBlogRoutes(router)
}
