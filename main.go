package main

import (
	"log"

	"github.com/kake87/CF_UserRegistrationService/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDatabase() {
	dsn := "host=localhost user=postgres password=yourpassword dbname=yourdb port=5432 sslmode=disable"

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Could not connect to the database &v", err)
	}

	log.Println("✅Succesful connect")
}

func main() {
	initDatabase()
	db.AutoMigrate(&models.User{})
}
