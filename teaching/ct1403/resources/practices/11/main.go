package main

import (
	"database/sql"
	"errors"
	"example/internal/httpserver"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/crypto/bcrypt"
)

func connectToDB() (*sql.DB, error) {
	dbUrl := os.Getenv("POSTGRESQL_URL")
	if dbUrl == "" {
		return nil, errors.New("POSTGRESQL_URL environment is not set!")
	}
	db, err := sql.Open("pgx", dbUrl)
	return db, err
}

func main() {
	db, err := connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if len(os.Args) < 2 {
		fmt.Println("Please enter an argument serve or create-user")
		return
	}
	switch os.Args[1] {
	case "serve":
		httpserver.ServeBlog(db)
	case "create-user":
		if len(os.Args) < 5 {
			fmt.Println("Please enter username, password and display_name")
			return
		}
		encryptedPass, err := bcrypt.GenerateFromPassword([]byte(os.Args[3]), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal(err)
		}
		_, err = db.Query(`INSERT INTO users (user_name, encrypted_password, display_name)
				VALUES ($1, $2, $3);`, os.Args[2], encryptedPass, os.Args[4])
		if err != nil {
			log.Fatal(err)
		}
	}
}
