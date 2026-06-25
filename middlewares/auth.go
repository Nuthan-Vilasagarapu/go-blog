package middlewares

import (
	"go-blog/constants"
	"go-blog/stores"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	session := stores.Session
	return func(ctx *gin.Context) {
		authdata := ctx.Request.Header.Get("Authorization")
		c_session := sessions.Default(ctx)
		var c_token string = ""
		if c_session.Get(constants.TokenKey) != nil {
			c_token = c_session.Get(constants.TokenKey).(string)
		}
		if authdata == "" && c_token == "" {
			ctx.JSON(http.StatusForbidden, gin.H{
				"message": "Not logged in",
			})
			ctx.Abort()
		}
		authdata_arr := strings.SplitN(authdata, "Bearer ", 2)
		if (len(authdata_arr) < 2 || strings.TrimSpace(authdata_arr[1]) == "") && c_token == "" {
			ctx.JSON(http.StatusForbidden, gin.H{
				"message": "Invalid Session",
			})
			ctx.Abort()
			return
		}
		var f_token string = ""
		if c_token != "" {
			f_token = c_token
		} else if len(authdata_arr) == 2 {
			f_token = authdata_arr[1]
		}
		user_id := session[f_token]
		if user_id == "" {
			ctx.JSON(http.StatusForbidden, gin.H{
				"message": "Invalid Session",
			})
			ctx.Abort()
			return
		}
		ctx.Set("user_id", user_id)
		ctx.Next()
	}
}
