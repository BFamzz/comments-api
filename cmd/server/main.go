package main

import (
	"fmt"
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

	fmt.Println("successfully created and pinged the database")

	return nil
}

func main() {
	fmt.Println("Comments API")
	if err := Run(); err != nil {
		fmt.Println(err)
	}
}
