package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Static("", "./www")
	router.Run("127.0.0.1:9000")
}
