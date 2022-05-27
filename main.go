package main

import (
	"log"
	"os"

	"github.com/carbondesinged/backstage-features-service/database"
	"github.com/carbondesinged/backstage-features-service/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	database.ConnectDb()
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
