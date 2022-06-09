package controllers

import (
	"context"

	"github.com/carbondesigned/backstage-features-service/config"
	"github.com/carbondesigned/backstage-features-service/database"
	"github.com/carbondesigned/backstage-features-service/models"
	"github.com/carbondesigned/backstage-features-service/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func GetImages(c *fiber.Ctx) error {
 var images []models.Image  

 	// if there is an error getting images from db
  if err := database.DB.Db.Find(&images).Error; err != nil {
    return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
      "success": false, 
      "error": err.Error(),
    })
  }

  return c.Status(fiber.StatusOK).JSON(fiber.Map{
    "success": true,
    "data": images,
  })
}

func CreateImage(c *fiber.Ctx) error {
  var image models.Image
  
  // able to read multipart/form-data form
	c.Accepts("multipart/form-data")
	c.Request().MultipartForm()

	if err := c.BodyParser(&image); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error trying to create image",
			"error":   err.Error(),
		})
	}

	token := c.Get("Authorization")
	claims, err := utils.ParseToken(token)

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
  
	newImage, err := c.FormFile("image")
	if err != nil {
	  return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	    "success": false,
	    "message": "Error when finding image in request",
	    "error": err.Error(),
	  })
	}

	// Uploading the image to a bucket.
	imageURL, err := config.UploadImage(context.TODO(), int(author.ID), newImage)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error trying to upload image",
			"error":   err.Error(),
		})
	}

	image.ImageURL = imageURL
	database.DB.Db.Create(&image)

  return c.Status(fiber.StatusOK).JSON(fiber.Map{
    "success": true,
    "data": image,
  })
}
