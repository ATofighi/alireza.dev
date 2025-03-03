package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/secret-area", gin.BasicAuth(gin.Accounts{
		"user1": "pass1",
		"user2": "pass2",
	}), func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)
		c.JSON(200, gin.H{
			"message": fmt.Sprintf("hii %s", user),
		})
	})
	r.Run(":9000")
}
