package main

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
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
type IBlogs struct {
	BlogName  string    `json:"blog_name"`
	Id        string    `json:"id"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type IUsers struct {
	Id           string    `json:"id"`
	UserName     string    `json:"username"`
	CreatedAt    time.Time `json:"created_at"`
	LastLoginAt  time.Time `json:"last_session_on"`
	passwordHash string
}

type DBBlogDataFmt struct {
	BlogName  string `yaml:"blog_name"`
	Content   string `yaml:"content"`
	Author    string `yaml:"author"`
	CreatedAt int64  `yaml:"created_at"`
	UpdatedAt int64  `yaml:"updated_at"`
}

type DBBlogFmt struct {
	Id   string        `yaml:"id"`
	Data DBBlogDataFmt `yaml:"data"`
}

type DBUserDataFmt struct {
	Username     string `yaml:"username"`
	PasswordHash string `yaml:"password_hash"`
	CreatedAt    int64  `yaml:"created_at"`
	LastLoginAt  int64  `yaml:"last_session_on"`
}

type DBUserFmt struct {
	Id   string        `yaml:"id"`
	Data DBUserDataFmt `yaml:"data"`
}

type DBFmt struct {
	Blogs []DBBlogFmt `yaml:"blogs"`
	Users []DBUserFmt `yaml:"users"`
}

func ReadFromDb(db *DBFmt, blogs *[]IBlogs) []IBlogs {
	for _, entry := range db.Blogs {
		*blogs = append(*blogs, IBlogs{
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

func ReadUsersFromDb(db *DBFmt, users *[]IUsers) []IUsers {
	for _, entry := range db.Users {
		*users = append(*users, IUsers{
			Id:        entry.Id,
			UserName:     entry.Data.Username,
			passwordHash: entry.Data.PasswordHash,
			CreatedAt: time.Unix(entry.Data.CreatedAt, 0),
			LastLoginAt:time.Unix(entry.Data.CreatedAt, 0),
		})
	}
	return *users
}

func WriteToDb(db *DBFmt, blogs []IBlogs) {
	db.Blogs = nil
	for _, entry := range blogs {
		db.Blogs = append(db.Blogs, DBBlogFmt{
			Id: entry.Id,
			Data: DBBlogDataFmt{
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

func WriteUsersToDb(db *DBFmt, users []IUsers) {
	db.Users = nil
	for _, entry := range users {
		db.Users = append(db.Users, DBUserFmt{
			Id: entry.Id,
			Data: DBUserDataFmt{
				Username  : entry.UserName,   
				PasswordHash : entry.passwordHash,
				CreatedAt :entry.CreatedAt.Unix(),   
				LastLoginAt :entry.LastLoginAt.Unix(),
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

func authMiddleware(session map[string]string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authdata := ctx.Request.Header.Get("Authorization")
		if authdata == "" {
			ctx.JSON(http.StatusForbidden, gin.H{
				"message": "Not logged in",
			})
			ctx.Abort()
		}
		authdata_arr := strings.Split(authdata, "Bearer ")
		fmt.Println(authdata_arr)
		if authdata_arr[1] == "" {
			ctx.JSON(http.StatusForbidden, gin.H{
				"message": "Invalid Session",
			})
			ctx.Abort()
		} else {
			user_id := session[authdata_arr[1]]
			if user_id == "" {
				ctx.JSON(http.StatusForbidden, gin.H{
					"message": "Invalid Session",
				})
				ctx.Abort()
			}
			ctx.Set("user_id", user_id)
		}
		ctx.Next()
	}
}

func main() {
	router := gin.Default()
	blogs := []IBlogs{}
	users := []IUsers{}
	session := make(map[string]string)

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
	users = ReadUsersFromDb(&db, &users)
	router.POST("/auth/register", func(ctx *gin.Context) {
		username := ctx.Request.FormValue("username")
		password_plain := ctx.Request.FormValue("password")
		hashBytes := md5.Sum([]byte(password_plain))
		password_hash := hex.EncodeToString(hashBytes[:])
		id := uuid.New()
		user := IUsers{
			Id:           id.String(),
			passwordHash: password_hash,
			UserName:     username,
			CreatedAt:    time.Now(),
		}
		for _, userexist := range users {
			if userexist.UserName == username {
				ctx.JSON(http.StatusConflict, gin.H{
					"message": "User already exists",
				})
				return
			}
		}
		users = append(users, user)
		WriteUsersToDb(&db,users)
		ctx.JSON(http.StatusCreated, gin.H{
			"message": "user created",
			"user":    user,
		})
	})

	router.POST("/auth/login", func(ctx *gin.Context) {
		username := ctx.Request.FormValue("username")
		password_plain := ctx.Request.FormValue("password")
		hashBytes := md5.Sum([]byte(password_plain))
		password_hash := hex.EncodeToString(hashBytes[:])

		for i, user := range users {
			if user.passwordHash == password_hash && user.UserName == username {
				token := make([]byte, 16)
				_, err := rand.Read(token)
				if err != nil {
					fmt.Println("Error in Rand")
				}
				token_str := hex.EncodeToString(token)
				session[token_str] = user.Id
				users[i].LastLoginAt = time.Now()
				WriteUsersToDb(&db,users)
				ctx.JSON(http.StatusAccepted, gin.H{
					"message": "Login successful",
					"token":   token_str,
				})
				return
			}
		}
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Login failed",
		})

	})

	router.GET("/users", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"users": users,
		})
	})

	router.POST("/blog", authMiddleware(session), func(ctx *gin.Context) {
		user_id, _ := ctx.Get("user_id")
		author_id := user_id.(string)
		blog_name := ctx.Request.FormValue("blog_name")
		content := ctx.Request.FormValue("content")
		id := uuid.New()
		blog := IBlogs{
			BlogName:  blog_name,
			Id:        id.String(),
			Content:   content,
			Author:    author_id,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		blogs = append(blogs, blog)
		WriteToDb(&db, blogs)
		ctx.JSON(http.StatusCreated, gin.H{
			"blog": blog,
		})
	})

	router.GET("/blog/:id", authMiddleware(session),func(ctx *gin.Context) {
		id := ctx.Params.ByName("id")
		var blogitem IBlogs
		for _, blog := range blogs {
			if blog.Id == id {
				blogitem = blog
			}
		}
		if blogitem == (IBlogs{}) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "blog not found",
			})
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Blog found",
				"blog":    blogitem,
			})
		}
	})

	router.GET("/search", authMiddleware(session),func(ctx *gin.Context) {
		author := ctx.Query("author")
		name := ctx.Query("name")
		content := ctx.Query("content")
		BlogFound := []IBlogs{}
		for _, blog := range blogs {
			if (name != "" && strings.Contains(blog.BlogName, name)) || (content != "" && strings.Contains(blog.Content, content)) || (author != "" && strings.Contains(blog.Author, author)) {
				BlogFound = append(BlogFound, blog)
			}
		}
		if len(BlogFound) == 0 {
			ctx.JSON(http.StatusNoContent, gin.H{
				"message": "No Blogs found with search",
			})
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Blogs matched with search",
				"Blogs":   BlogFound,
			})
		}
	})

	router.PUT("/blog/:id",authMiddleware(session), func(ctx *gin.Context) {
		id := ctx.Params.ByName("id")
		name := ctx.Request.FormValue("name")
		content := ctx.Request.FormValue("content")
		author := ctx.Request.FormValue("author")

		for i, blog := range blogs {
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
				blogs[i] = blog
				WriteToDb(&db, blogs)
				ctx.JSON(http.StatusOK, gin.H{
					"message": "Blog updated successfully",
					"Blogs":   blogs,
				})
				return
			}
		}
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Blog not found, please check blogs from below to be updated",
			"blogs":   blogs,
		})

	})

	router.DELETE("/blog/:id", authMiddleware(session), func(ctx *gin.Context) {
		id := ctx.Params.ByName("id")
		var blogitem IBlogs
		var pos int
		for i, blog := range blogs {
			if blog.Id == id {
				blogitem = blog
				pos = i
			}
		}
		if blogitem == (IBlogs{}) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "blog not found",
			})
		} else {
			blogs = slices.Delete(blogs, pos, pos+1)
			WriteToDb(&db, blogs)
			ctx.JSON(http.StatusNoContent, gin.H{
				"message": "Blog deleted succesfully",
				"blogs":   blogs,
			})
		}
	})

	router.DELETE("/blogs", authMiddleware(session), func(ctx *gin.Context) {
		blogs = []IBlogs{}
		WriteToDb(&db, blogs)
		ctx.JSON(http.StatusNoContent, gin.H{
			"message": "All Blogs deleted succesfully, there are no blogs",
			"blogs":   blogs,
		})
	})

	router.GET("/blogs", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"blogs": blogs,
		})
	})

	router.SetTrustedProxies(nil)
	router.Run(":8000")
}
