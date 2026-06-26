package main

import (
	"go-blog/routers"

	"github.com/gin-gonic/gin"
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

func main() {
	router := gin.Default()

	routers.LoadRoutes(router)

	router.SetTrustedProxies(nil)
	router.Run(":8000")
}
