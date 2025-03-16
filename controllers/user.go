package controllers

import (
	"closeloop/config"
	"closeloop/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RegisterUser(c *fiber.Ctx) error {
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if err := user.Validate(); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// add hashed password into db

	result := config.DB.Create(user)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": result.Error.Error()})
	}

	return c.Status(201).JSON(
		fiber.Map{
			"status": "OK",
			"data":   "User created successfully",
		},
	)
}

func LoginUser(c *fiber.Ctx) error {
	// Struct for receiving login data
	var input struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	// Fetch user from database using either email or username
	var user models.User
	query := config.DB

	if input.Email != "" {
		query = query.Where("email = ?", input.Email)
	} else if input.Username != "" {
		query = query.Where("username = ?", input.Username)
	} else {
		return c.Status(400).JSON(fiber.Map{"error": "Email or username is required"})
	}

	result := query.First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}

	// Compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expires in 1 minute
	})

	tokenString, err := token.SignedString(config.JWTSecret)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not generate token"})
	}

	// Return token to user
	return c.Status(200).JSON(fiber.Map{
		"status":  "OK",
		"message": "User logged in successfully",
		"token":   tokenString,
	})
}

func GetUser(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{"status": "OK", "data": "user"})
}
