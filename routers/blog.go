package routers

import (
	"go-blog/controllers"
	"go-blog/middlewares"

	"github.com/gin-gonic/gin"
)

func LoadBlogRoutes(router *gin.Engine) {
	router.Use(middlewares.AuthMiddleware())
	// API
	router.POST("/blog", controllers.CreateBlog)
	router.GET("/blog/:id", controllers.GetBlogByID)
	router.GET("/search", controllers.SearchBlogs)
	router.PUT("/blog/:id", controllers.UpdateBlog)
	router.DELETE("/blog/:id", controllers.DeleteBlog)
	router.DELETE("/blogs", controllers.DeleteBlogs)
	router.GET("/blogs", controllers.ListBlogs)

	// HTML
	router.GET("/views/blog/:id", controllers.ViewBlog)
	router.GET("/views/blogs", controllers.ViewBlogs)
	router.GET("/add/blog", controllers.AddBlogPage)
}
