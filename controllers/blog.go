package controllers

import (
	"go-blog/interfaces"
	"go-blog/repository"
	"go-blog/stores"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateBlog(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	authorID := userID.(string)
	blogName := ctx.Request.FormValue("blog_name")
	content := ctx.Request.FormValue("content")
	id := uuid.New()
	blog := interfaces.IBlogs{
		BlogName:  blogName,
		Id:        id.String(),
		Content:   content,
		Author:    authorID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	stores.Blogs = append(stores.Blogs, blog)
	repository.WriteToDb(&stores.DB, stores.Blogs)
	ctx.JSON(http.StatusCreated, gin.H{"blog": blog})
}

func GetBlogByID(ctx *gin.Context) {
	id := ctx.Params.ByName("id")
	for _, blog := range stores.Blogs {
		if blog.Id == id {
			ctx.JSON(http.StatusOK, gin.H{"message": "Blog found", "blog": blog})
			return
		}
	}
	ctx.JSON(http.StatusNotFound, gin.H{"message": "blog not found"})
}

func SearchBlogs(ctx *gin.Context) {
	author := ctx.Query("author")
	name := ctx.Query("name")
	content := ctx.Query("content")
	blogFound := []interfaces.IBlogs{}
	for _, blog := range stores.Blogs {
		if (name != "" && strings.Contains(blog.BlogName, name)) || (content != "" && strings.Contains(blog.Content, content)) || (author != "" && strings.Contains(blog.Author, author)) {
			blogFound = append(blogFound, blog)
		}
	}
	if len(blogFound) == 0 {
		ctx.JSON(http.StatusNoContent, gin.H{"message": "No Blogs found with search"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Blogs matched with search", "Blogs": blogFound})
}

func UpdateBlog(ctx *gin.Context) {
	id := ctx.Params.ByName("id")
	name := ctx.Request.FormValue("name")
	content := ctx.Request.FormValue("content")
	author := ctx.Request.FormValue("author")

	for i, blog := range stores.Blogs {
		if blog.Id == id {
			if name != "" {
				blog.BlogName = name
			}
			if content != "" {
				blog.Content = content
			}
			if author != "" {
				blog.Author = author
			}
			blog.UpdatedAt = time.Now()
			stores.Blogs[i] = blog
			repository.WriteToDb(&stores.DB, stores.Blogs)
			ctx.JSON(http.StatusOK, gin.H{"message": "Blog updated successfully", "Blogs": stores.Blogs})
			return
		}
	}
	ctx.JSON(http.StatusNotFound, gin.H{"message": "Blog not found", "blogs": stores.Blogs})
}

func DeleteBlog(ctx *gin.Context) {
	id := ctx.Params.ByName("id")
	for i, blog := range stores.Blogs {
		if blog.Id == id {
			stores.Blogs = slices.Delete(stores.Blogs, i, i+1)
			repository.WriteToDb(&stores.DB, stores.Blogs)
			ctx.JSON(http.StatusNoContent, gin.H{"message": "Blog deleted succesfully", "blogs": stores.Blogs})
			return
		}
	}
	ctx.JSON(http.StatusNotFound, gin.H{"message": "blog not found"})
}

func DeleteBlogs(ctx *gin.Context) {
	stores.Blogs = []interfaces.IBlogs{}
	repository.WriteToDb(&stores.DB, stores.Blogs)
	ctx.JSON(http.StatusNoContent, gin.H{"message": "All Blogs deleted succesfully", "blogs": stores.Blogs})
}

func ListBlogs(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"blogs": stores.Blogs})
}

func AddBlogPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "add-blog.html", gin.H{})
}

func ViewBlog(ctx *gin.Context) {
	id := ctx.Params.ByName("id")
	authors := make(map[string]string, 10)
	for _, author := range stores.Users {
		authors[author.Id] = author.UserName
	}
	for _, blog := range stores.Blogs {
		viewBlog := interfaces.IViewBlog{
			Name:        blog.BlogName,
			Content:     blog.Content,
			Author:      authors[blog.Author],
			PublishedAt: time.Unix(blog.CreatedAt.Unix(), 0).UTC().String(),
		}
		if id == blog.Id {
			ctx.HTML(http.StatusOK, "view-blog.html", gin.H{
				"title": "Blog Website",
				"Blog":  viewBlog,
			})
			return
		}
	}
	ctx.JSON(http.StatusNotFound, gin.H{"message": "Not Found!!!"})
}

func ViewBlogs(ctx *gin.Context) {
	blogs := interfaces.IViewBlogs{}
	authors := make(map[string]string, 10)
	for _, author := range stores.Users {
		authors[author.Id] = author.UserName
	}
	for _, blog := range stores.Blogs {
		viewBlog := interfaces.IViewBlog{
			Name:        blog.BlogName,
			Content:     blog.Content,
			Author:      authors[blog.Author],
			PublishedAt: time.Unix(blog.CreatedAt.Unix(), 0).UTC().String(),
			Id:          blog.Id,
		}
		blogs.Blogs = append(blogs.Blogs, viewBlog)
	}
	ctx.HTML(http.StatusOK, "list-blog.html", gin.H{
		"title": "Blog Website",
		"Blogs": blogs.Blogs,
	})
}
