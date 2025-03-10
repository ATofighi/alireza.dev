package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.Static("/", "./www")
	r.Run("127.0.0.1:3000")
}
