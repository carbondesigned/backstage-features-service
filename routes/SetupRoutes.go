package routes

import (
	"github.com/carbondesigned/backstage-features-service/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Setup() *fiber.App {
	database.ConnectDb()

	app := fiber.New()
	app.Use(cors.New())

	SetupRoutes(app)

	return app
}

func SetupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "Success! 🚀",
		})
	})
	api := app.Group("/api")

	AuthRoutes(api.Group("/auth"))
	PostRoutes(api.Group("/posts"))
	AlbumRoutes(api.Group("/albums"))
	ImageRoutes(api.Group("/images"))
}
