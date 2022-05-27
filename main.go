package main

import (
	"log"

	"github.com/digitalocean/sample-golang/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	// app.Use(cors.New())

	routes.SetupRoutes(app)

	log.Println("Server is running on port: 8080")
	err := app.Listen(":80")
	if err != nil {
		log.Fatal(err)
	}
}
