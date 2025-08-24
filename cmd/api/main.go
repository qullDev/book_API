package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/qullDev/book_API/internal/cache"
	"github.com/qullDev/book_API/internal/config"
	"github.com/qullDev/book_API/internal/db"
	"github.com/qullDev/book_API/internal/domain/book"
	"github.com/qullDev/book_API/internal/domain/category"
	"github.com/qullDev/book_API/internal/domain/user"
	"github.com/qullDev/book_API/internal/http/router"
	appauth "github.com/qullDev/book_API/internal/pkg/auth"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/qullDev/book_API/docs" // swagger docs
)

// @title Book API
// @version 1.0
// @description RESTful service for managing books and categories with JWT authentication
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	// set mode gin (debug/release)
	if cfg.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Connect to database
	dbConn, err := db.Connect(cfg)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	// Connect to redis
	rdb, err := cache.Connect(cfg)
	if err != nil {
		log.Fatal("Error connecting to redis:", err)
	}
	ts := appauth.NewTokenStore(rdb)

	// AutoMigrate
	if err := dbConn.AutoMigrate(&user.User{}, &category.Category{}, &book.Book{}); err != nil {
		log.Fatal("Error migrating database:", err)
	}
	log.Println("✅ Database migrated")

	// SEED user
	var count int64
	dbConn.Model(&user.User{}).Count(&count)
	if count == 0 {
		hashed, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		dbConn.Create(&user.User{
			Username: "admin",
			Password: string(hashed),
		})
	}
	r := router.New(dbConn, cfg, ts)
	log.Println("Server is running on port:", cfg.AppPort)

	// Update to use PORT env var from Railway
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
