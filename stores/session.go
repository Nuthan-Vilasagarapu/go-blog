// session := make(map[string]string)
package stores

import "go-blog/interfaces"

var (
	Session = make(map[string]string)
	Blogs   = []interfaces.IBlogs{}
	Users   = []interfaces.IUsers{}
)
