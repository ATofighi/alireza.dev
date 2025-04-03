package httpserver

import (
	"database/sql"
	"example/internal/db/models"
	"example/internal/db/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ServeBlog(db *sql.DB) {

	postsRepository := repositories.NewPostRepository(db)
	const title = "بلاگ تستی"

	r := gin.Default()
	r.LoadHTMLGlob("./templates/*")
	r.Static("/static", "./static")
	r.GET("/", func(ctx *gin.Context) {
		posts, err := postsRepository.GetAll()
		if err != nil {
			ctx.Error(err)
			return
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
		post := models.Post{
			Title:    ctx.PostForm("title"),
			Content:  ctx.PostForm("content"),
			AuthorId: 1,
		}
		err := postsRepository.Insert(&post)
		if err != nil {
			ctx.Error(err)
			return
		}
		ctx.Redirect(http.StatusFound, "/")
	})
	r.Run("127.0.0.1:3000")
}
