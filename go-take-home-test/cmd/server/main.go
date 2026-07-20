package main

import (
	"log"
	"os"

	"go-take-home-test/internal/app"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	e := app.New()
	log.Printf("Server is running on http://localhost:%s", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatal(err)
	}
}
