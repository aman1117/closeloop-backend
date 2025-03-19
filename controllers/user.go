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

type APIResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
}

func RegisterUser(c *fiber.Ctx) error {
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(400).JSON(APIResponse{Status: "Error", Error: err.Error()})
	}

	if err := user.Validate(); err != nil {
		return c.Status(400).JSON(APIResponse{Status: "Error", Error: err.Error()})
	}

	if result := config.DB.Create(user); result.Error != nil {
		return c.Status(500).JSON(APIResponse{Status: "Error", Error: result.Error.Error()})
	}

	return c.Status(201).JSON(APIResponse{Status: "OK", Data: "User created successfully"})
}

func LoginUser(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(APIResponse{Status: "Error", Error: "Invalid request payload"})
	}

	var user models.User
	query := config.DB

	if input.Email != "" {
		query = query.Where("email = ?", input.Email)
	} else if input.Username != "" {
		query = query.Where("username = ?", input.Username)
	} else {
		return c.Status(400).JSON(APIResponse{Status: "Error", Error: "Email or username is required"})
	}

	if err := query.First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(401).JSON(APIResponse{Status: "Error", Error: "Invalid credentials"})
		}
		return c.Status(500).JSON(APIResponse{Status: "Error", Error: "Database error"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.Status(401).JSON(APIResponse{Status: "Error", Error: "Invalid credentials"})
	}

	if len(config.JWTSecret) == 0 {
		return c.Status(500).JSON(APIResponse{Status: "Error", Error: "JWT secret is not configured"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		return c.Status(500).JSON(APIResponse{Status: "Error", Error: "Could not generate token"})
	}

	return c.Status(200).JSON(APIResponse{
		Status: "OK",
		Data:   fiber.Map{"token": tokenString},
	})
}

func GetUser(c *fiber.Ctx) error {
	var user models.User
	username, ok := c.Locals("username").(string)
	if !ok {
		return c.Status(400).JSON(APIResponse{Status: "Error", Error: "Invalid token payload"})
	}

	if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(APIResponse{Status: "Error", Error: "User not found"})
		}
		return c.Status(500).JSON(APIResponse{Status: "Error", Error: "Database error"})
	}

	userResponse := struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Avatar   string `json:"avatar"`
	}{
		ID:       user.ID.String(),
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Avatar:   user.Avatar,
	}

	return c.Status(200).JSON(APIResponse{Status: "OK", Data: userResponse})
}
