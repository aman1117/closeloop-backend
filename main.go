package main

import (
	"closeloop/config"
	"closeloop/models"
	"closeloop/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.ConnectDB()

	app := fiber.New()
	config.DB.AutoMigrate(&models.User{})
	routes.SetupRoutes(app)

	app.Listen(":3000")
}
