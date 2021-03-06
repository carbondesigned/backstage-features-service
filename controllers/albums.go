package controllers

import (
	"context"

	"github.com/carbondesigned/backstage-features-service/config"
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

	// is the usera an author
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

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Not authorized",
			"error":   err.Error(),
		})
	}
	// we process the image and upload it to a bucket
	cover, err := c.FormFile("cover")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error trying with cover image",
			"error":   err.Error(),
		})
	}
	// Uploading the image to a bucket.
	coverURL, err := config.UploadImage(context.TODO(), int(author.ID), cover)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error trying to upload image",
			"error":   err.Error(),
		})
	}

	// set the image URL from bucket to post's cover.
	album.Cover = coverURL

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

func UploadToAlbum(c *fiber.Ctx) error {
	var album models.Album
	var newAlbum models.Album

	c.Accepts("multipart/form-data")

	albumSlug := c.Params("id")

	if err := c.BodyParser(&newAlbum); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error trying to create album",
			"error":   err.Error(),
		})
	}

	// find the album
	if err := database.DB.Db.Where("slug = ?", albumSlug).First(&album).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error trying to get album",
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

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Not authorized",
			"error":   err.Error(),
		})
	}

	// take a single image and add it to an array of the album's images
	image, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error trying with image",
			"error":   err.Error(),
		})
	}
	// Uploading the image to a bucket.
	imageURL, err := config.UploadImage(context.TODO(), int(author.ID), image)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error trying to upload image",
			"error":   err.Error(),
		})
	}

	// add the image to the album's images array
	album.Images = append(album.Images, imageURL)

	// update the album
	if err := database.DB.Db.Save(&album).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error trying to update album",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    album,
	})
}

func GetAlbum(c *fiber.Ctx) error {
	// get album by it's slug
	var album models.Album
	albumSlug := c.Params("id")

	if err := database.DB.Db.Where("slug = ?", albumSlug).First(&album).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error trying to get album",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    album,
	})
}

func EditAlbum(c *fiber.Ctx) error {
	var album models.Album
	var newAlbum models.Album
	c.Accepts("multipart/form-data")

	albumSlug := c.Params("id")

	if err := c.BodyParser(&newAlbum); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error trying to create album",
			"error":   err.Error(),
		})
	}

	// find the album
	if err := database.DB.Db.Where("slug = ?", albumSlug).First(&album).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error trying to get album",
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

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Not authorized",
			"error":   err.Error(),
		})
	}

	// we process the image and upload it to a bucket
	cover, err := c.FormFile("cover")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error trying with cover image",
			"error":   err.Error(),
		})
	}

	if album.Cover != newAlbum.Cover {
		coverURL, err := config.UploadImage(context.TODO(), int(author.ID), cover)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Error trying to upload image",
				"error":   err.Error(),
			})
		}

		newAlbum.CoverURL = coverURL
	}

	album.Slug = utils.GenerateSlugFromTitle(newAlbum.Title)

	if err := database.DB.Db.Model(&album).Updates(newAlbum).Error; err != nil {
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

func DeleteAlbum(c *fiber.Ctx) error {
	var album models.Album

	// get post slug from url
	slug := c.Params("id")

	if err := database.DB.Db.Where("slug = ?", slug).First(&album).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Album not found",
			"error":   err.Error(),
		})
	}

	// get user from token
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

	// if the user doesn't exist, they can't delete a post (because they are not an author)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized to delete album",
			"error":   err.Error(),
		})
	}

	database.DB.Db.Delete(&album)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    album,
	})
}
