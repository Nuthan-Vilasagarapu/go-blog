package controllers

import (
	"go-blog/interfaces"
	"go-blog/repository"
	"go-blog/stores"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateBlog(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	authorID := userID.(string)
	blogName := ctx.Request.FormValue("blog_name")
	content := ctx.Request.FormValue("content")

	blog := repository.CreateBlog(blogName, content, authorID)

	ctx.JSON(http.StatusCreated, gin.H{"blog": blog})
}

func GetBlogByID(ctx *gin.Context) {
	id := ctx.Params.ByName("id")
	blog := repository.GetBlogById(id)
	if blog != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "Blog found", "blog": blog})
	} else {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "blog not found"})
	}
}

func SearchBlogs(ctx *gin.Context) {
	name := ctx.Query("name")
	content := ctx.Query("content")
	blogFound := repository.SearchBlogs(name, content)
	if len(blogFound) == 0 {
		ctx.JSON(http.StatusNoContent, gin.H{"message": "No Blogs found with search"})
		ctx.Abort()
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Blogs matched with search", "Blogs": blogFound})
}

func UpdateBlog(ctx *gin.Context) {
	id := ctx.Params.ByName("id")
	name := ctx.Request.FormValue("name")
	content := ctx.Request.FormValue("content")

	status := repository.UpdateBlogByID(id, interfaces.IBlog{
		BlogName: name,
		Content:  content,
	})
	if status != true {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Blog not found"})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": "Blog updated!"})
	}
}

func DeleteBlog(ctx *gin.Context) {
	id := ctx.Params.ByName("id")
	status := repository.DeleteBlog(id)
	if status {
		ctx.JSON(http.StatusNoContent, gin.H{"message": "Blog deleted succesfully", "blogs": stores.Blogs})
	} else {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "blog not found"})
	}
}

func DeleteBlogs(ctx *gin.Context) {
	_ = repository.DeleteBlogs()
	ctx.JSON(http.StatusNoContent, gin.H{"message": "All Blogs deleted succesfully", "blogs": stores.Blogs})
}

func ListBlogs(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"blogs": repository.GetAllBlogs()})
}

func AddBlogPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "add-blog.html", gin.H{})
}

func ViewBlog(ctx *gin.Context) {
	id := ctx.Params.ByName("id")

	blog := repository.GetBlogById(id)
	if blog == nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "Blog not found!"})
	}

	viewBlog := interfaces.IViewBlog{
		Id:          blog.Id,
		Name:        blog.BlogName,
		Content:     blog.Content,
		Author:      repository.GetUserById(blog.Author).UserName,
		PublishedAt: time.Unix(blog.CreatedAt.Unix(), 0).UTC().String(),
	}
	ctx.HTML(http.StatusOK, "view-blog.html", gin.H{
		"title": "Blog Website",
		"Blog":  viewBlog,
	})
}

func ViewBlogs(ctx *gin.Context) {
	var blogs []interfaces.IViewBlog
	for _, blog := range repository.GetAllBlogs() {
		blogs = append(blogs, interfaces.IViewBlog{
			Id:          blog.Id,
			Name:        blog.BlogName,
			Content:     blog.Content,
			Author:      repository.GetUserById(blog.Author).UserName,
			PublishedAt: time.Unix(blog.CreatedAt.Unix(), 0).UTC().String(),
		})
	}
	ctx.HTML(http.StatusOK, "list-blog.html", gin.H{
		"title": "Blog Website",
		"Blogs": blogs,
	})
}
