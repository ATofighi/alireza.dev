package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Post struct {
	PublishedDate string
	Title         string
	Content       string
}

func main() {
	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS posts (
		id integer primary key,
		title text,
		content text,
		date text
	)`)
	if err != nil {
		log.Fatal(err)
	}

	const title = "بلاگ تستی"

	r := gin.Default()
	r.LoadHTMLGlob("./templates/*")
	r.Static("/static", "./static")
	r.GET("/", func(ctx *gin.Context) {
		posts := make([]Post, 0)
		rows, err := db.Query("SELECT title, content, date from posts order by id desc")
		if err != nil {
			log.Fatal(err)
		}
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
		_, err := db.Exec("INSERT INTO posts (title, content, date) values (?, ?, ?)",
			ctx.PostForm("title"),
			ctx.PostForm("content"),
			time.Now().Format(time.RFC850),
		)
		if err != nil {
			log.Fatal(err)
		}
		ctx.Redirect(http.StatusFound, "/")
	})
	r.Run("127.0.0.1:3000")
}
