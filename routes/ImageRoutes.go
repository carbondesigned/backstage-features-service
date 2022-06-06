package routes

import (
	"github.com/carbondesigned/backstage-features-service/controllers"
	"github.com/gofiber/fiber/v2"
)


func ImageRoutes(route fiber.Router) {
  route.Get("/", controllers.GetImages)
  route.Post("/create", controllers.CreateImage)
}
