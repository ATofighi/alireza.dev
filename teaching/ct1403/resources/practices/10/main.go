package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Post struct {
	PublishedDate string
	Title         string
	Content       string
}

func main() {
	db, err := sql.Open("pgx", "postgres://ct1403:123456@localhost:5432/blog")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	const title = "بلاگ تستی"

	r := gin.Default()
	r.LoadHTMLGlob("./templates/*")
	r.Static("/static", "./static")
	r.GET("/", func(ctx *gin.Context) {
		post_id := ctx.Query("post_id")
		posts := make([]Post, 0)
		rows, err := db.Query("SELECT title, content, created_at from posts WHERE id=$1 order by id desc", post_id)
		if err != nil {
			log.Println(err)
		} else {
			defer rows.Close()
			for rows.Next() {
				post := Post{}
				rows.Scan(&post.Title, &post.Content, &post.PublishedDate)
				posts = append(posts, post)
			}
			err = rows.Err()
			if err != nil {
				ctx.Error(err)
			}
		}

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
		_, err := db.Exec("INSERT INTO posts (author_id, title, content) values (1, $1, $2)",
			ctx.PostForm("title"),
			ctx.PostForm("content"),
		)
		if err != nil {
			log.Println(err)
		}
		ctx.Redirect(http.StatusFound, "/")
	})
	r.Run("127.0.0.1:3000")
}
