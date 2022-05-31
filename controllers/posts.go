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

// We create a post, we check if the user is an author, we upload the image to a bucket, we set the
// coverURL to the post, we generate a slug from the title, and we create the post
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

	// we process the image and upload it to a bucket
	cover := post.Cover

	// Uploading the image to a bucket.
	coverURL, err := config.UploadImage(context.TODO(), int(author.ID), cover)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error trying to upload image",
			"error":   err.Error(),
		})
	}
	post.Cover = coverURL
	post.Slug = utils.GenerateSlugFromTitle(post.Title)
	database.DB.Db.Create(&post)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    post,
	})
}

// We get the post from the database, get the user from the token, and if the user exists,
// we update the post
func EditPost(c *fiber.Ctx) error {
	var post models.Post
	var newPost models.Post

	// get post slug from url
	slug := c.Params("id")

	if err := c.BodyParser(&newPost); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error trying to edit post",
			"error":   err.Error(),
		})
	}

	// get post from database
	err := database.DB.Db.Where("slug = ?", slug).First(&post).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Post not found",
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

	// if the user doesn't exist, they can't edit a post (because they are not an author)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized to edit post",
			"error":   err.Error(),
		})
	}

	// If the title of the post is different from the new title, we generate a new slug from the new title
	if post.Title != newPost.Title {
		newPost.Slug = utils.GenerateSlugFromTitle(newPost.Title)
	}
	database.DB.Db.Model(&post).Updates(newPost)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    post,
	})
}

// It gets a post by its slug and increments its views by one
func GetPost(c *fiber.Ctx) error {
	var post models.Post
	slug := c.Params("id")

	err := database.DB.Db.Where("slug = ?", slug).First(&post).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Post not found",
			"error":   err.Error(),
		})
	}

	database.DB.Db.Model(&post).Update("views", post.Views+1)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    post,
	})
}

// We get the post slug from the URL, check if the post exists, get the user from the token, check if
// the user exists, and if they do, delete the post
func DeletePost(c *fiber.Ctx) error {
	var post models.Post

	// get post slug from url
	slug := c.Params("id")

	if err := database.DB.Db.Where("slug = ?", slug).First(&post).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Post not found",
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
			"message": "Unauthorized to delete post",
			"error":   err.Error(),
		})
	}

	database.DB.Db.Delete(&post)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    post,
	})
}

func LikePost(c *fiber.Ctx) error {
	var post models.Post
	slug := c.Params("id")

	database.DB.Db.Where(
		"slug = ?",
		slug,
	).Model(&post).Update("likes", post.Likes+1)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    post,
	})
}
