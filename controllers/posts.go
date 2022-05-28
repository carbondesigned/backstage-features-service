package controllers

import (
	"context"

	"github.com/carbondesigned/backstage-features-service/config"
	"github.com/carbondesigned/backstage-features-service/database"
	"github.com/carbondesigned/backstage-features-service/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func GetPosts(c *fiber.Ctx) error {
	var posts []models.Post
	err := database.DB.Db.Find(&posts).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    posts,
	})
}

// We're creating a new post, but only if the user is authorized to do so
func CreatePost(c *fiber.Ctx) error {
	var post models.Post

	if err := c.BodyParser(&post); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error trying to create post",
			"error":   err.Error(),
		})
	}

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
	author := models.Author{}
	err = database.DB.Db.Where("id = ?", id).First(&author).Error

	// if the user doesn't exist, they can't create a post (because they are not an author)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized to create post",
			"error":   err.Error(),
		})
	}

	// we process the image and upload it to a bucket
	cover := post.Cover
	coverURL, err := config.UploadImage(context.TODO(), int(author.ID), cover)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error trying to upload image",
			"error":   err.Error(),
		})
	}
	// we set the coverURL to the post
	post.CoverURL = coverURL

	database.DB.Db.Create(&post)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    post,
	})
}
