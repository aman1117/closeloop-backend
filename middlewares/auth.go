package middleware

import (
	"closeloop/config"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type APIResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
}

func AuthMiddleware(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(401).JSON(APIResponse{Status: "error", Error: "Missing token"})
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(config.JWTSecret), nil // Ensure this is a byte slice
	})

	if err != nil || !token.Valid {
		return c.Status(401).JSON(APIResponse{Status: "error", Error: "Invalid token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(401).JSON(APIResponse{Status: "error", Error: "Invalid token claims"})
	}

	if exp, ok := claims["exp"].(float64); ok {
		if int64(exp) < time.Now().Unix() {
			return c.Status(401).JSON(APIResponse{Status: "error", Error: "Token expired"})
		}
	}
	c.Locals("username", claims["username"])
	return c.Next()
}
