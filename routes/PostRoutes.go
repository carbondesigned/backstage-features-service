package routes

import (
	"github.com/carbondesigned/backstage-features-service/controllers"
	"github.com/gofiber/fiber/v2"
)

func PostRoutes(route fiber.Router) {
	route.Get("/", controllers.GetPosts)
	route.Post("/create", controllers.CreatePost)
	route.Put("/:id", controllers.EditPost)
	route.Get("/:id", controllers.GetPost)
	route.Delete("/:id", controllers.DeletePost)
}
