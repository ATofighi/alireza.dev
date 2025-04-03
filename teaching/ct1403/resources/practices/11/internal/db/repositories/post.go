package repositories

import (
	"database/sql"
	"example/internal/db/models"
)

type PostRepository interface {
	GetAll() ([]models.Post, error)
	Insert(*models.Post) error
}

type IPostRepository struct {
	db *sql.DB
}

func (r *IPostRepository) GetAll() ([]models.Post, error) {
	rows, err := r.db.Query("SELECT title, content, created_at, author_id from posts order by id desc")
	posts := make([]models.Post, 0)
	if err != nil {
		return posts, err
	}
	defer rows.Close()
	for rows.Next() {
		post := models.Post{}
		rows.Scan(&post.Title, &post.Content, &post.PublishedDate, &post.AuthorId)
		posts = append(posts, post)
	}
	err = rows.Err()
	return posts, err
}

func (r *IPostRepository) Insert(post *models.Post) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
		INSERT INTO posts (title, content, author_id)
		VALUES($1, $2, $3)
		`, post.Title, post.Content, post.AuthorId)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
		UPDATE users set post_count=post_count+1
		WHERE id=$1
		`, post.AuthorId)
	if err != nil {
		return err
	}
	err = tx.Commit()
	return err
}

func NewPostRepository(db *sql.DB) PostRepository {
	return &IPostRepository{
		db: db,
	}
}
