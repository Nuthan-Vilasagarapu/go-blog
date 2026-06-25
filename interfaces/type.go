package interfaces

import "time"

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
	PasswordHash string
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

type IViewBlog struct {
	Id, Name, Content, Author, PublishedAt string
}

type IViewBlogs struct {
	Blogs []IViewBlog
}
