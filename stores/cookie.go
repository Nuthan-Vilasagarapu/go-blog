package stores

import (
	"go-blog/constants"

	"github.com/gin-gonic/gin"
)

func SetAuthCookie(ctx *gin.Context, token string) {
	ctx.SetCookie(constants.CookieName, token, 3600, "/", "", false, true)
}

func ReadAuthCookie(ctx *gin.Context) (string, error) {
	return ctx.Cookie(constants.CookieName)
}

func ClearAuthCookie(ctx *gin.Context) {
	ctx.SetCookie(constants.CookieName, "", -1, "/", "", false, true)
}
