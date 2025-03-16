package middleware

import (
	"closeloop/config"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Missing token"})
	}

	//GO CONCEPT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return config.JWTSecret, nil
	})

	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
	}

	// Check if the token is valid
	if !token.Valid {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
	}

	// Extract claims if needed
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid token claims"})
	}

	// Check expiration
	if exp, ok := claims["exp"].(float64); ok {
		if int64(exp) < time.Now().Unix() {
			return c.Status(401).JSON(fiber.Map{"error": "Token expired"})
		}
	}

	return c.Next()
}
