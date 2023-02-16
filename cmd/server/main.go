package main

import (
	"context"
	"fmt"
	"github.com/BFamzz/comments-api/internal/comment"
	"github.com/BFamzz/comments-api/internal/db"
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
	fmt.Println(commentService.GetComment(context.Background(), "276c1dcc-b800-4801-a870-380d891dbc6a"))

	return nil
}

func main() {
	fmt.Println("Comments API")
	if err := Run(); err != nil {
		fmt.Println(err)
	}
}
