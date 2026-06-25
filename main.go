package main

import (
	"go-blog/middlewares"
	"go-blog/routers"
	"go-blog/stores"

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

	stores.LoadDB()
	middlewares.LoadMiddlewares(router)
	routers.LoadRoutes(router)

	router.SetTrustedProxies(nil)
	router.Run(":8000")
}
