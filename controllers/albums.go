package controllers

import (
	"context"
	"sync"

	"github.com/carbondesigned/backstage-features-service/config"
	database "github.com/carbondesigned/backstage-features-service/database"
	models "github.com/carbondesigned/backstage-features-service/models"
	"github.com/carbondesigned/backstage-features-service/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func CreateAlbum(c *fiber.Ctx) error {
	var album models.Album
	var wg sync.WaitGroup

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

	if album.Cover != "" {
		// we process the image and upload it to a bucket
		cover := album.Cover
		wg.Add(1)
		coverURL, err := config.UploadImage(&wg, context.TODO(), int(author.ID), cover)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Error trying to upload image",
				"error":   err.Error(),
			})
		}
		// we set the coverURL to the post
		album.CoverURL = coverURL
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

func UploadToAlbum(c *fiber.Ctx) error {
	var album models.Album
	var newAlbum models.Album
	var wg sync.WaitGroup

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
	images := newAlbum.RootImages

	for _, image := range images {
		imageURL, err := config.UploadImage(&wg, context.TODO(), int(author.ID), image)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Error trying to upload image",
				"error":   err.Error(),
			})
		}
		// we set the coverURL to the post
		newAlbum.Images = append(newAlbum.Images, imageURL)
	}

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

func EditAlbum(c *fiber.Ctx) error {
	var album models.Album
	var newAlbum models.Album
	var wg sync.WaitGroup

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

	if album.Cover != newAlbum.Cover {
		coverURL, err := config.UploadImage(&wg, context.TODO(), int(author.ID), newAlbum.Cover)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Error trying to upload image",
				"error":   err.Error(),
			})
		}
		// we set the coverURL to the post
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
