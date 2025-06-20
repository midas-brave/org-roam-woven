package main

import (
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to org-roam-woven!",
		})
	})

	return r
}

func main() {
	r := setupRouter()

	r.Run(":18080")
}
