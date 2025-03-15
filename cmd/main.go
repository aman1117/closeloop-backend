package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		log.Fatal("Missing DB_PASSWORD environment variable")
	}
	dsn := fmt.Sprintf("host=db.svnqvehncqwhpnavjykh.supabase.co user=postgres password=%s dbname=postgres port=5432 sslmode=require", dbPassword)

	_, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Database connected successfully!")
	app := fiber.New()

	api := app.Group("/api")
	api.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "OK"})
	})

	app.Listen(":3000")
}
