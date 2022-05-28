package controllers

import (
	database "github.com/carbondesigned/backstage-features-service/database"
	models "github.com/carbondesigned/backstage-features-service/models"
	"github.com/carbondesigned/backstage-features-service/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func CreateAlbum(c *fiber.Ctx) error {
	var album models.Album

	if err := c.BodyParser(&album); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error trying to create album",
			"error":   err.Error(),
		})
	}

	// is the user an author
	token := c.Get("Authorization")
	claims, err := jwt.Parse(
		token,
		func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		},
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Invalid token",
			"error":   err.Error(),
		})
	}

	id := claims.Claims.(jwt.MapClaims)["id"]

	err = database.DB.Db.Where("id = ?", id).First(&models.Author{}).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Not authorized",
			"error":   err.Error(),
		})
	}

	album.Slug = utils.GenerateSlugFromTitle(album.Title)
	if err := database.DB.Db.Create(&album).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error trying to create album",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    album,
	})
}

func GetAlbums(c *fiber.Ctx) error {
	var albums []models.Album

	if err := database.DB.Db.Find(&albums).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error trying to get albums",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    albums,
	})
}
