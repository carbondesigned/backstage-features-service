package main

import (
	"log"
	"os"

	"github.com/digitalocean/sample-golang/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	// app.Use(cors.New())

	routes.SetupRoutes(app)

	// Get the PORT from  env
	port := os.Getenv("PORT")

	// Verify if heroku provided the port or not
	if os.Getenv("PORT") == "" {
		port = "3000"
	}

	// Start server on http://${heroku-url}:${port}
	log.Fatal(app.Listen(":" + port))
}
