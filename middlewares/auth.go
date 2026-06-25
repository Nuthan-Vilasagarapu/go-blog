package middlewares

import (
	"go-blog/stores"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Read token from the Authorization header first, then fall back to the cookie.
		authHeader := ctx.GetHeader("Authorization")
		token := ""

		if authHeader != "" {
			parts := strings.SplitN(authHeader, "Bearer ", 2)
			if len(parts) == 2 && strings.TrimSpace(parts[1]) != "" {
				token = strings.TrimSpace(parts[1])
			}
		}

		if token == "" {
			var err error
			token, err = stores.ReadAuthCookie(ctx)
			if err != nil || strings.TrimSpace(token) == "" {
				ctx.JSON(http.StatusForbidden, gin.H{"message": "Not logged in"})
				ctx.Abort()
				return
			}
		}

		// Validate the token against the in-memory session map.
		userID := stores.Session[token]
		if userID == "" {
			ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid Session"})
			ctx.Abort()
			return
		}

		// Attach the authenticated user ID to the request context.
		ctx.Set("user_id", userID)
		ctx.Next()
	}
}
