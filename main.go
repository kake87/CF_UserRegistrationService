package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kake87/CF_UserRegistrationService/models" // Импорт моделей из GitHub
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDatabase() {
	// Подключение к базе данных PostgreSQL
	dsn := "host=localhost user=postgres password=yourpassword dbname=yourdb port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	log.Println("Connected to the database")
}

func main() {
	// Инициализация базы данных
	initDatabase()

	// Миграция модели User
	db.AutoMigrate(&models.User{})

	router := gin.Default()

	// Маршрут для регистрации пользователя
	router.POST("/register", func(c *gin.Context) {
		var input models.User

		// Парсинг JSON из тела запроса
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Хэширование пароля
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Error hashing password:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		// Создание пользователя в базе данных
		user := models.User{Name: input.Name, Email: input.Email, Password: string(hashedPassword)}
		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		// Успешный ответ
		c.JSON(http.StatusOK, gin.H{"message": "User registered successfully", "user_id": user.ID})
	})

	router.Run(":8082") // Микросервис работает на порту 8081
}
