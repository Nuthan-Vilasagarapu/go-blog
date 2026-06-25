package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HomePage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"title":   "Blog Website",
		"message": "Wellcome to my Blogs",
	})
}
