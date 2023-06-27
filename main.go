package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)
	dsn := "host=postgres user=root password=root dbname=root sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	sqlDB, _ := db.DB()

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// // SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	// sqlDB.SetConnMaxLifetime(time.Hour)

	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&User{})

	app := fiber.New()

	app.Post("/", func(c *fiber.Ctx) error {
		var userInput UserInput
		c.BodyParser(&userInput)
		db.Create(&User{
			Name:      userInput.Name,
			UserName:  userInput.UserName,
			Email:     userInput.Email,
			Bio:       userInput.Bio,
			BirthDate: userInput.BirthDate,
			Gender:    userInput.Gender,
			Avatar:    userInput.Avatar,
			Header:    userInput.Header,
			Password:  userInput.Password,
		})

		return c.JSON(UserInput{
			Name:      userInput.Name,
			UserName:  userInput.UserName,
			Email:     userInput.Email,
			Bio:       userInput.Bio,
			BirthDate: userInput.BirthDate,
			Gender:    userInput.Gender,
			Avatar:    userInput.Avatar,
			Header:    userInput.Header,
			Password:  userInput.Password,
		})
	})

	app.Get("/", func(c *fiber.Ctx) error {
		var users []User

		var UsersCount int64
		db.Find(&User{}).Count(&UsersCount)
		rand := rand.Intn(int(UsersCount) - 50)
		db.Offset(rand).Limit(50).Find(&users)

		return c.JSON(users)
	})
	fmt.Println("Server Running on 5000")
	app.Listen(":5000")
}

type UserInput struct {
	Name      string `json:"name"`
	UserName  string `json:"userName"`
	Email     string `json:"email"`
	Bio       string `json:"bio"`
	BirthDate string `json:"birthDate"`
	Gender    string `json:"gender"`
	Avatar    string `json:"avatar"`
	Header    string `json:"header"`
	Password  string `json:"password"`
}

type User struct {
	gorm.Model
	Name      string `json:"name"`
	UserName  string `json:"userName"`
	Email     string `json:"email"`
	Bio       string `json:"bio"`
	BirthDate string `json:"birthDate"`
	Gender    string `json:"gender"`
	Avatar    string `json:"avatar"`
	Header    string `json:"header"`
	Password  string `json:"password"`
}
