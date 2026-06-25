package stores

import (
	"fmt"
	"go-blog/interfaces"
	"go-blog/repository"
	"os"

	"github.com/goccy/go-yaml"
)

var (
	DB = interfaces.DBFmt{}
)

func LoadDB() {
	blogDB, err := os.ReadFile("blogDB.yaml")
	if err != nil {
		_, createErr := os.Create("blogDB.yaml")
		if createErr != nil {
			fmt.Println("failed to create blogDB.yaml:", createErr)
		}
	}
	err = yaml.Unmarshal(blogDB, &DB)
	if err != nil {
		fmt.Println("Yaml Error:", err)
	}
	Blogs = repository.ReadFromDb(&DB, &Blogs)
	Users = repository.ReadUsersFromDb(&DB, &Users)
}
