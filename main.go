package main

import (
	"log"
	"os"

	"github.com/carbondesigned/backstage-features-service/routes"
)

func main() {
	app := routes.Setup()
	// Get the PORT from  env
	port := os.Getenv("PORT")

	// Verify if port is available
	if os.Getenv("PORT") == "" {
		port = "80"
	}

	log.Fatal(app.Listen(":" + port))
}
