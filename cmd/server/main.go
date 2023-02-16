package main

import (
	"fmt"
	"github.com/BFamzz/comments-api/internal/comment"
	"github.com/BFamzz/comments-api/internal/db"
	transportHttp "github.com/BFamzz/comments-api/internal/transport/http"
)

func Run() error {
	fmt.Println("Starting up comments api application")

	database, err := db.NewDatabase()
	if err != nil {
		fmt.Println("failed to connect to the database")
		return err
	}
	if err := database.MigrateDB(); err != nil {
		fmt.Println("failed to migrate database")
		return err
	}

	commentService := comment.NewService(database)

	httpHandler := transportHttp.NewHandler(commentService)
	if err := httpHandler.Serve(); err != nil {
		return err
	}

	return nil
}

func main() {
	fmt.Println("Comments API")
	if err := Run(); err != nil {
		fmt.Println(err)
	}
}
