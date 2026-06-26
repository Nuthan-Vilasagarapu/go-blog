package interfaces

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
	UpdatedAt    int    `yaml:"last_updated_at"`
}

type DBUserFmt struct {
	Id   string        `yaml:"id"`
	Data DBUserDataFmt `yaml:"data"`
}

type DBFmt struct {
	Blogs []DBBlogFmt `yaml:"blogs"`
	Users []DBUserFmt `yaml:"users"`
}
