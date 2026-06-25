// SearchBlog

package repository

import (
	"go-blog/interfaces"
	"go-blog/stores"
	"slices"
	"strings"

	"time"

	"github.com/google/uuid"
)

func GetAllBlogs() []interfaces.IBlog {
	return stores.Blogs
}

func GetBlogById(id string) *interfaces.IBlog {
	for _, blog := range GetAllBlogs() {
		viewBlog := interfaces.IBlog{
			BlogName:  blog.BlogName,
			Content:   blog.Content,
			Author:    blog.Author,
			Id:        blog.Id,
			CreatedAt: blog.CreatedAt,
			UpdatedAt: blog.UpdatedAt,
		}
		if id == viewBlog.Id {
			return &viewBlog
		}
	}
	return nil
}

func CreateBlog(
	name, content, author string,
) interfaces.IBlog {
	id := uuid.New()
	blog := interfaces.IBlog{
		BlogName:  name,
		Id:        id.String(),
		Content:   content,
		Author:    author,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	stores.Blogs = append(stores.Blogs, blog)
	WriteToDb(&stores.DB, stores.Blogs)
	return blog
}

func DeleteBlog(
	id string,
) bool {
	for i, blog := range GetAllBlogs() {
		if blog.Id == id {
			stores.Blogs = slices.Delete(stores.Blogs, i, i+1)
			WriteToDb(&stores.DB, stores.Blogs)
			return true
		}
	}
	return false
}

func UpdateBlogByID(
	id string,
	new_blog interfaces.IBlog,
) bool {
	for i, blog := range GetAllBlogs() {
		if blog.Id == id {
			if new_blog.BlogName != "" {
				blog.BlogName = new_blog.BlogName
			}
			if new_blog.Content != "" {
				blog.Content = new_blog.Content
			}
			blog.UpdatedAt = time.Now()
			stores.Blogs[i] = blog
			WriteToDb(&stores.DB, stores.Blogs)
			return true
		}
	}
	return false
}

func DeleteBlogs() bool {
	stores.Blogs = []interfaces.IBlog{}
	WriteToDb(&stores.DB, stores.Blogs)
	return true
}

func GetBlogsByAuthorId(
	author_id string,
) []interfaces.IBlog {
	authorBlogs := []interfaces.IBlog{}

	for _, blog := range GetAllBlogs() {
		if blog.Author == author_id {
			authorBlogs = append(authorBlogs, blog)
		}
	}

	return authorBlogs
}

func SearchBlogs(
	name_snippet, content_snippet string,
) []interfaces.IBlog {
	blogFound := []interfaces.IBlog{}
	for _, blog := range GetAllBlogs() {
		if (name_snippet != "" && strings.Contains(blog.BlogName, name_snippet)) || (content_snippet != "" && strings.Contains(blog.Content, content_snippet)) {
			blogFound = append(blogFound, blog)
		}
	}
	return blogFound
}
