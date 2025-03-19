package config

import (
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var JWTSecret = []byte(os.Getenv("AMAN_SECRET_KEY"))

func ConnectDB() {

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		log.Fatal("Missing DB_PASSWORD environment variable")
	}
	dsn := fmt.Sprintf("host=closeloopdb.postgres.database.azure.com user=postgres password=%s dbname=postgres port=5432 ", dbPassword)

	//DOCS
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("🚀 Connected to the database!")
	DB = db
}
