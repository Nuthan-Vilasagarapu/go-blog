package stores

import (
	"go-blog/constants"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
)

func CookieStore() cookie.Store {
	store := cookie.NewStore([]byte(constants.SecretKey))
	store.Options(
		sessions.Options{
			Path:     "/",
			Secure:   false,
			HttpOnly: true,
			MaxAge:   3600,
		},
	)

	return store
}
