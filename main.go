package main

import (
	"github.com/gin-gonic/gin"
)

func startServer() error {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to org-roam-woven!",
		})
	})

	println("Starting server on :18080")
	return r.Run(":18080")
}

func main() {
	if err := startServer(); err != nil {
		println("Error starting server:", err.Error())
	}
}

