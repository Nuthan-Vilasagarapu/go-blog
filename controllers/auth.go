package controllers

import (
	"go-blog/interfaces"
	"go-blog/repository"
	"go-blog/stores"
	"go-blog/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterUser(ctx *gin.Context) {
	username := ctx.Request.FormValue("username")
	passwordPlain := ctx.Request.FormValue("password")

	passwordHash := utils.HashPassword(passwordPlain)
	user, err := repository.CreateUser(username, passwordHash)
	if err != nil {
		ctx.JSON(http.StatusCreated, gin.H{
			"message": "registration failed",
			"error":   err.Error(),
		})
		ctx.Abort()
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "user created",
		"user":    user,
	})
}

func LoginUser(ctx *gin.Context) {
	username := ctx.Request.FormValue("username")
	passwordPlain := ctx.Request.FormValue("password")

	user := repository.GetUserByUsername(username)
	if user == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		ctx.Abort()
	}

	if !utils.Comparehash(passwordPlain, user.Data.PasswordHash) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Invalid Password"})
		ctx.Abort()
	}

	token, err := utils.GenerateSessionToken()
	if err != nil {
		ctx.JSON(http.StatusCreated, gin.H{
			"message": "login failed",
			"error":   err.Error(),
		})
		ctx.Abort()
	}
	stores.SetSession(*token, user.Id)

	success := repository.UpdateUser(user.Id, interfaces.DBUserDataFmt{
		LastLoginAt: time.Now().Unix(),
	})

	if success {
		stores.SetAuthCookie(ctx, *token)
		ctx.JSON(http.StatusAccepted, gin.H{
			"message": "Login successful",
			"token":   token,
		})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{
			"message": "Login Failed",
		})
	}
}

func GetUsers(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"users": repository.GetUsers(),
	})
}

func LoginPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "auth.html", gin.H{})
}
