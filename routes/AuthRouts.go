package routes

import (
	"github.com/digitalocean/sample-golang/controllers"
	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(route fiber.Router) {
	route.Post("signin", controllers.Signin)
}
