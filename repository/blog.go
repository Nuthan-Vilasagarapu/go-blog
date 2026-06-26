// SearchBlog

package repository

import (
	"go-blog/interfaces"
	"slices"
	"strings"

	"time"

	"github.com/google/uuid"
)

func GetAllBlogs() []interfaces.DBBlogFmt {
	return ReadFromDb()
}

func GetBlogById(id string) *interfaces.DBBlogFmt {
	for _, blog := range GetAllBlogs() {
		viewBlog := interfaces.DBBlogFmt{
			Id: blog.Id,
			Data: interfaces.DBBlogDataFmt{
				BlogName:  blog.Data.BlogName,
				Content:   blog.Data.Content,
				Author:    blog.Data.Author,
				CreatedAt: blog.Data.CreatedAt,
				UpdatedAt: blog.Data.UpdatedAt,
			},
		}
		if id == viewBlog.Id {
			return &viewBlog
		}
	}
	return nil
}

func CreateBlog(
	name, content, author string,
) interfaces.DBBlogFmt {
	blogs := GetAllBlogs()
	id := uuid.New()
	blog := interfaces.DBBlogFmt{
		Id: id.String(),
		Data: interfaces.DBBlogDataFmt{
			BlogName:  name,
			Content:   content,
			Author:    author,
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		},
	}
	blogs = append(blogs, blog)
	WriteToDb(blogs)
	return blog
}

func DeleteBlog(
	id string,
) bool {
	blogs := GetAllBlogs()
	for i, blog := range blogs {
		if blog.Id == id {
			blogs = slices.Delete(blogs, i, i+1)
			WriteToDb(blogs)
			return true
		}
	}
	return false
}

func UpdateBlogByID(
	id string,
	new_blog interfaces.DBBlogDataFmt,
) bool {
	blogs := GetAllBlogs()
	for i, blog := range blogs {
		if blog.Id == id {
			if new_blog.BlogName != "" {
				blog.Data.BlogName = new_blog.BlogName
			}
			if new_blog.Content != "" {
				blog.Data.Content = new_blog.Content
			}
			blog.Data.UpdatedAt = time.Now().Unix()
			blogs[i] = blog
			WriteToDb(blogs)
			return true
		}
	}
	return false
}

func DeleteBlogs() bool {
	Blogs := []interfaces.DBBlogFmt{}
	WriteToDb(Blogs)
	return true
}

func GetBlogsByAuthorId(
	author_id string,
) []interfaces.DBBlogFmt {
	authorBlogs := []interfaces.DBBlogFmt{}
	for _, blog := range GetAllBlogs() {
		if blog.Data.Author == author_id {
			authorBlogs = append(authorBlogs, blog)
		}
	}
	return authorBlogs
}

func SearchBlogs(
	name_snippet, content_snippet string,
) []interfaces.DBBlogFmt {
	blogFound := []interfaces.DBBlogFmt{}
	for _, blog := range GetAllBlogs() {
		if (name_snippet != "" && strings.Contains(blog.Data.BlogName, name_snippet)) || (content_snippet != "" && strings.Contains(blog.Data.Content, content_snippet)) {
			blogFound = append(blogFound, blog)
		}
	}
	return blogFound
}
