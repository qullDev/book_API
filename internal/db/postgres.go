package db

import (
	"fmt"
	"log"

	"github.com/qullDev/book_API/internal/config"
	"github.com/qullDev/book_API/internal/domain/book"
	"github.com/qullDev/book_API/internal/domain/category"
	"github.com/qullDev/book_API/internal/domain/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s  password=%s sslmode=%s", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBName, cfg.DBPassword, cfg.DBSSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database:", err)
		return nil, err
	}

	// Auto migrate schema
	if err := db.AutoMigrate(&category.Category{}, &book.Book{}, &user.User{}); err != nil {
		log.Println("AutoMigrate warning:", err)
	}

	return db, nil

}
