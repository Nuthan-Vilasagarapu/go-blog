package middlewares

import "github.com/gin-gonic/gin"

func LoadMiddlewares(router *gin.Engine) {
	router.Use(cookieMiddleware())
}
