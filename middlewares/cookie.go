// sessions.Sessions(sessionName, store)

package middlewares

import (
	"go-blog/constants"
	"go-blog/stores"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func cookieMiddleware() gin.HandlerFunc {
	store := stores.CookieStore()
	return sessions.Sessions(constants.SessionName, store)
}
