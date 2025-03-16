package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var JWTSecret = []byte(os.Getenv("AMAN_SECRET_KEY"))

func ConnectDB() {
	//FIX
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		log.Fatal("Missing DB_PASSWORD environment variable")
	}
	dsn := fmt.Sprintf("host=db.svnqvehncqwhpnavjykh.supabase.co user=postgres password=%s dbname=postgres port=5432 sslmode=require", dbPassword)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("🚀 Connected to the database!")
	DB = db
}
