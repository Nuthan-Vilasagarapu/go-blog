package repository

import (
	"fmt"
	"go-blog/interfaces"
	"os"
	"time"

	"github.com/goccy/go-yaml"
)

func ReadFromDb(db *interfaces.DBFmt, blogs *[]interfaces.IBlogs) []interfaces.IBlogs {
	for _, entry := range db.Blogs {
		*blogs = append(*blogs, interfaces.IBlogs{
			BlogName:  entry.Data.BlogName,
			Id:        entry.Id,
			Content:   entry.Data.Content,
			Author:    entry.Data.Author,
			CreatedAt: time.Unix(entry.Data.CreatedAt, 0),
			UpdatedAt: time.Unix(entry.Data.UpdatedAt, 0),
		})
	}

	return *blogs
}

func ReadUsersFromDb(db *interfaces.DBFmt, users *[]interfaces.IUsers) []interfaces.IUsers {
	for _, entry := range db.Users {
		*users = append(*users, interfaces.IUsers{
			Id:           entry.Id,
			UserName:     entry.Data.Username,
			PasswordHash: entry.Data.PasswordHash,
			CreatedAt:    time.Unix(entry.Data.CreatedAt, 0),
			LastLoginAt:  time.Unix(entry.Data.CreatedAt, 0),
		})
	}
	return *users
}

func WriteToDb(db *interfaces.DBFmt, blogs []interfaces.IBlogs) {
	db.Blogs = nil
	for _, entry := range blogs {
		db.Blogs = append(db.Blogs, interfaces.DBBlogFmt{
			Id: entry.Id,
			Data: interfaces.DBBlogDataFmt{
				BlogName:  entry.BlogName,
				Content:   entry.Content,
				Author:    entry.Author,
				CreatedAt: entry.CreatedAt.Unix(),
				UpdatedAt: entry.UpdatedAt.Unix(),
			},
		})
	}
	data, err := yaml.Marshal(db)
	if err != nil {
		fmt.Println("failed to marshal db:", err)
		return
	}
	os.WriteFile("blogDB.yaml", data, 0644)
}

func WriteUsersToDb(db *interfaces.DBFmt, users []interfaces.IUsers) {
	db.Users = nil
	for _, entry := range users {
		db.Users = append(db.Users, interfaces.DBUserFmt{
			Id: entry.Id,
			Data: interfaces.DBUserDataFmt{
				Username:     entry.UserName,
				PasswordHash: entry.PasswordHash,
				CreatedAt:    entry.CreatedAt.Unix(),
				LastLoginAt:  entry.LastLoginAt.Unix(),
			},
		})
	}
	data, err := yaml.Marshal(db)
	if err != nil {
		fmt.Println("failed to marshal db:", err)
		return
	}
	os.WriteFile("blogDB.yaml", data, 0644)
}
