package controllers

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"go-blog/constants"
	"go-blog/interfaces"
	"go-blog/repository"
	"go-blog/stores"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterUser(ctx *gin.Context) {
	username := ctx.Request.FormValue("username")
	passwordPlain := ctx.Request.FormValue("password")
	hashBytes := md5.Sum([]byte(passwordPlain))
	passwordHash := hex.EncodeToString(hashBytes[:])
	id := uuid.New()
	user := interfaces.IUsers{
		Id:           id.String(),
		PasswordHash: passwordHash,
		UserName:     username,
		CreatedAt:    time.Now(),
	}
	for _, userExist := range stores.Users {
		if userExist.UserName == username {
			ctx.JSON(http.StatusConflict, gin.H{
				"message": "User already exists",
			})
			ctx.Abort()
			return
		}
	}
	stores.Users = append(stores.Users, user)
	repository.WriteUsersToDb(&stores.DB, stores.Users)
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "user created",
		"user":    user,
	})
}

func LoginUser(ctx *gin.Context) {
	cSession := sessions.Default(ctx)
	username := ctx.Request.FormValue("username")
	passwordPlain := ctx.Request.FormValue("password")
	hashBytes := md5.Sum([]byte(passwordPlain))
	passwordHash := hex.EncodeToString(hashBytes[:])

	for i, user := range stores.Users {
		if user.PasswordHash == passwordHash && user.UserName == username {
			token := make([]byte, 16)
			if _, err := rand.Read(token); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
				return
			}
			tokenStr := hex.EncodeToString(token)
			stores.Session[tokenStr] = user.Id
			stores.Users[i].LastLoginAt = time.Now()
			repository.WriteUsersToDb(&stores.DB, stores.Users)
			cSession.Set(constants.TokenKey, tokenStr)
			if err := cSession.Save(); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
				return
			}
			ctx.JSON(http.StatusAccepted, gin.H{
				"message": "Login successful",
				"token":   tokenStr,
			})
			return
		}
	}
	ctx.JSON(http.StatusUnauthorized, gin.H{
		"message": "Login failed",
	})
}

func GetUsers(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"users": stores.Users,
	})
}

func LoginPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "auth.html", gin.H{})
}
