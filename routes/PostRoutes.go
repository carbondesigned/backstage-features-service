package routes

import (
	"github.com/carbondesigned/backstage-features-service/controllers"
	"github.com/gofiber/fiber/v2"
)

func PostRoutes(route fiber.Router) {
	route.Get("/", controllers.GetPosts)
	route.Post("/create", controllers.CreatePost)
}
