package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/qullDev/book_API/internal/cache"
	"github.com/qullDev/book_API/internal/config"
	"github.com/qullDev/book_API/internal/db"
	"github.com/qullDev/book_API/internal/http/router"
)

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
	if _, err := db.Connect(cfg); err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	// Connect to redis
	if _, err := cache.Connect(cfg); err != nil {
		log.Fatal("Error connecting to redis:", err)
	}

	r := router.New()
	log.Println("Server is running on port:", cfg.AppPort)

	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
