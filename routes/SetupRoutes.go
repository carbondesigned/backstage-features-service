package routes

import (
	"github.com/gofiber/fiber/v2"
)

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
}
