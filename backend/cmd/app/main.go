package main

import (
	"log"

	"secret-santa-backend/internal/app"
)

func main() {
	app := app.New()

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
