package routes

import (
	"github.com/carbondesinged/backstage-features-service/controllers"
	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(route fiber.Router) {
	route.Post("/createuser", controllers.CreateUser)
	route.Post("/signin", controllers.Signin)
}
