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

	log.Fatal(app.Listen(":3001"))
}
