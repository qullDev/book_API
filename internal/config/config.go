package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort         string
	DBHost          string
	DBPort          string
	DBUser          string
	DBPassword      string
	DBName          string
	DBSSLMode       string
	RedisAddr       string
	RedisPassword   string
	RedisDB         int
	JWTSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	Env             string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	// default Redis DB jadi 1 (bukan 0)
	redisDB, _ := strconv.Atoi(getenv("REDIS_DB", "1"))

	at, err := time.ParseDuration(getenv("ACCESS_TOKEN_TTL", "15m"))
	if err != nil {
		at = 15 * time.Minute
	}
	rt, err := time.ParseDuration(getenv("REFRESH_TOKEN_TTL", "168h"))
	if err != nil {
		rt = 168 * time.Hour
	}

	return &Config{
		AppPort:         getenv("APP_PORT", "8080"),
		DBHost:          getenv("DB_HOST", "localhost"),
		DBPort:          getenv("DB_PORT", "5432"),
		DBUser:          getenv("DB_USER", "postgres"),
		DBPassword:      getenv("DB_PASSWORD", "12345"),
		DBName:          getenv("DB_NAME", "books_db"),
		DBSSLMode:       getenv("DB_SSLMODE", "disable"),
		RedisAddr:       getenv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:   getenv("REDIS_PASSWORD", ""),
		RedisDB:         redisDB,
		JWTSecret:       getenv("JWT_SECRET", "devsecret"),
		AccessTokenTTL:  at,
		RefreshTokenTTL: rt,
		Env:             getenv("ENV", "dev"),
	}, nil
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
