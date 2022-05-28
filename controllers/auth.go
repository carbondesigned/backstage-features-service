package controllers

import (
	"time"

	database "github.com/carbondesigned/backstage-features-service/database"
	models "github.com/carbondesigned/backstage-features-service/models"
	"github.com/carbondesigned/backstage-features-service/utils"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
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

// It creates a new author
func CreateAuthor(c *fiber.Ctx) error {
	author := new(models.Author)
	if err := c.BodyParser(author); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	foundAuthor := models.Author{}
	// if user already exists
	if err := database.DB.Db.Where("email = ?", foundAuthor).First(&author).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"success": false,
			"error":   "User already exists",
		})
	}
	// if password is less than 6 characters
	if len(author.Password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Password must be at least 6 characters",
		})
	}
	hashedPassword, err := utils.HashPassword(author.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	author.Password = hashedPassword
	database.DB.Db.Create(&author)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    author,
	})
}

// It takes the email and password from the request body, checks if the user exists, checks if the
// password is correct, creates a JWT token, and sends it back to the user
func Signin(c *fiber.Ctx) error {
	author := new(models.Author)
	if err := c.BodyParser(&author); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}
	email := author.Email
	password := author.Password
	var foundAuthor models.Author

	// check if user exists
	if err := database.DB.Db.Where("email = ?", email).First(&foundAuthor).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "User not found",
			"error":   err.Error(),
		})
	}
	// check if password is correct
	if !utils.CheckPasswordHash(password, foundAuthor.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Incorrect password",
		})
	}
	// Create the Claims
	claims := jwt.MapClaims{
		"id":    foundAuthor.ID,
		"email": email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:  "token",
		Value: t,
	})
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": email + " signed in",
		"token":   t,
	})
}
