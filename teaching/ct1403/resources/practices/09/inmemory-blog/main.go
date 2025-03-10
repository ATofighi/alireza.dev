package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Post struct {
	PublishedDate string
	Title         string
	Content       string
}

func main() {
	const title = "بلاگ تستی"
	posts := []Post{
		Post{
			Title:         "پست ۱",
			Content:       "محتوای تستییی",
			PublishedDate: "دیروز",
		},
		Post{
			Title:         "پست ۲",
			Content:       "محتوای تستییی",
			PublishedDate: "دو روز پیش",
		},
	}

	r := gin.Default()
	r.LoadHTMLGlob("./templates/*")
	r.Static("/static", "./static")
	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"title": title,
			"posts": posts,
		})
	})
	adminRouter := r.Group("/admin", gin.BasicAuth(gin.Accounts{
		"user": "pass",
	}))
	adminRouter.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "admin.html", gin.H{
			"title": title,
		})
	})
	adminRouter.POST("/", func(ctx *gin.Context) {
		posts = append(posts, Post{
			Title:         ctx.PostForm("title"),
			Content:       ctx.PostForm("content"),
			PublishedDate: time.Now().Format(time.RFC850),
		})
		ctx.Redirect(http.StatusFound, "/")
	})
	r.Run("127.0.0.1:3000")
}
