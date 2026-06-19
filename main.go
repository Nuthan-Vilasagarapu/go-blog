package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
	//"net/http"
	"os"
	//"slices"
	//"strings"
	"time"
)

/*
YAML content specification
blogs:

	id:
		blog_name:
		content
		author
		created_at
		updated_at
*/
type CBlogs struct {
	Blog_name string    `json:"blog_name"`
	Id        string    `json:"id"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DBBlogDataFmt struct {
	Blog_name string    `yaml:"blog_name"`
	Content   string    `yaml:"content"`
	Author    string    `yaml:"author"`
	CreatedAt int64 	`yaml:"created_at"`
	UpdatedAt int64 	`yaml:"updated_at"`
}

type DBBlogFmt struct {
	Id   string        `yaml:"id"`
	Data DBBlogDataFmt `yaml:"data"`
}

type DBFmt struct {
	Blogs []DBBlogFmt `yaml:"blogs"`
}

func ReadFromDb(db *DBFmt, blogs *[]CBlogs) []CBlogs {
	for _, entry := range db.Blogs {
		*blogs = append(*blogs, CBlogs{
			Blog_name: entry.Data.Blog_name,
			Id:        entry.Id,
			Content:   entry.Data.Content,    
			Author:    entry.Data.Author,
			CreatedAt: time.Unix(entry.Data.CreatedAt, 0),
			UpdatedAt: time.Unix(entry.Data.UpdatedAt, 0),
		})
	}

	return *blogs
}

func WriteToDb(db *DBFmt, todos []CBlogs) {
	for _, entry := range todos {
		db.Blogs = append(db.Blogs, DBBlogFmt{
			Id: entry.Id,
			Data: DBBlogDataFmt{
				Blog_name: entry.Blog_name,
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

func main(){
	router := gin.Default()
	blogs := []CBlogs{}
	db := DBFmt{}
	blogDB, err := os.ReadFile("blogDB.yaml")
	if err != nil {
		_, createErr := os.Create("blogDB.yaml")
		if createErr != nil {
			fmt.Println("failed to create blogDB.yaml:", createErr)
		}
	}
	err = yaml.Unmarshal(blogDB, &db)
	if err != nil {
		fmt.Println("Yaml Error:", err)
	}
	blogs = ReadFromDb(&db, &blogs)


	router.POST("/creating", func(ctx *gin.Context) {
		blog_name := ctx.Request.FormValue("blog_name")
		content   := ctx.Request.FormValue("content")
		author    := ctx.Request.FormValue("author")
		id := uuid.New()

		blog := CBlogs{
			Blog_name: blog_name,   
			Id:        id.String(),   
			Content:   content,    
			Author:    author,  
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		blogs = append(blogs, blog)
		WriteToDb(&db, blogs)
		ctx.JSON(200, gin.H{
			"blog": blog,
		})
	})

	router.GET("/create", func(ctx *gin.Context) {
		ctx.File("index.html")
	})

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"blogs": blogs,
		})
	})

	router.Run(":8000")

}