package controllers

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	// "github.com/golang-jwt/jwt/v4"
)

func AuthRequired() func(c *fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		SigningKey: "secrete",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		},
	})
}
func Signin(c *fiber.Ctx) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"success": false,
	})
}
