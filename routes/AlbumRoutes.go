package routes

import (
	"github.com/carbondesigned/backstage-features-service/controllers"
	"github.com/gofiber/fiber/v2"
)

func AlbumRoutes(route fiber.Router) {
	route.Post("/create", controllers.CreateAlbum)
	route.Get("/", controllers.GetAlbums)
}