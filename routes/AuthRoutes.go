package routes

import (
	"github.com/carbondesigned/backstage-features-service/controllers"
	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(route fiber.Router) {
	route.Post("/createuser", controllers.CreateAuthor)
	route.Post("/signin", controllers.Signin)
	route.Get("/", controllers.GetAuthors)
	route.Get("/:id", controllers.GetAuthor)
}
