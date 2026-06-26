package repository

import (
	"fmt"
	"go-blog/interfaces"
	"os"

	"github.com/goccy/go-yaml"
)

func ReadFromDb() []interfaces.DBBlogFmt {
	DB := interfaces.DBFmt{}

	blogDB, err := os.ReadFile("blogDB.yaml")
	if err != nil {
		_, createErr := os.Create("blogDB.yaml")
		if createErr != nil {
			fmt.Println("failed to create blogDB.yaml:", createErr)
		}
	}
	err = yaml.Unmarshal(blogDB, DB)
	if err != nil {
		fmt.Println("Yaml Error:", err)
	}

	return DB.Blogs
}

func ReadUsersFromDb() []interfaces.DBUserFmt {
	DB := interfaces.DBFmt{}

	blogDB, err := os.ReadFile("blogDB.yaml")
	if err != nil {
		_, createErr := os.Create("blogDB.yaml")
		if createErr != nil {
			fmt.Println("failed to create blogDB.yaml:", createErr)
		}
	}
	err = yaml.Unmarshal(blogDB, DB)
	if err != nil {
		fmt.Println("Yaml Error:", err)
	}

	return DB.Users
}

func WriteToDb(blogs []interfaces.DBBlogFmt) {
	DB := interfaces.DBFmt{
		Users: ReadUsersFromDb(),
	}

	for _, entry := range blogs {
		DB.Blogs = append(DB.Blogs, interfaces.DBBlogFmt{
			Id:   entry.Id,
			Data: entry.Data,
		})
	}
	data, err := yaml.Marshal(DB)
	if err != nil {
		fmt.Println("failed to marshal db:", err)
		return
	}
	os.WriteFile("blogDB.yaml", data, 0644)
}

func WriteUsersToDb(users []interfaces.DBUserFmt) {
	DB := interfaces.DBFmt{
		Blogs: ReadFromDb(),
	}

	for _, entry := range users {
		DB.Users = append(DB.Users, interfaces.DBUserFmt{
			Id:   entry.Id,
			Data: entry.Data,
		})
	}
	data, err := yaml.Marshal(DB)
	if err != nil {
		fmt.Println("failed to marshal db:", err)
		return
	}
	os.WriteFile("blogDB.yaml", data, 0644)
}
