package main

import "fmt"

func Run() error {
	fmt.Println("Starting up comments api application")
	return nil
}

func main() {
	fmt.Println("Comments API")
	if err := Run(); err != nil {
		fmt.Println(err)
	}
}
