package main

import (
	"log"
	"os"

	"github.com/digitalocean/sample-golang/db"
	"github.com/digitalocean/sample-golang/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	db.ConnectDb()
	app := fiber.New()
	app.Use(cors.New())

	routes.SetupRoutes(app)

	// Get the PORT from  env
	port := os.Getenv("PORT")

	// Verify if port is available
	if os.Getenv("PORT") == "" {
		port = "80"
	}

	log.Fatal(app.Listen(":" + port))
}
