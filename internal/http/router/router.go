package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qullDev/book_API/internal/config"
	"github.com/qullDev/book_API/internal/http/handlers"
	"github.com/qullDev/book_API/internal/http/middleware"
	appauth "github.com/qullDev/book_API/internal/pkg/auth"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func New(db *gorm.DB, cfg *config.Config, ts *appauth.TokenStore) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// Swagger route - pastikan ini ada di atas route lainnya
	r.GET("/api/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// endpoint health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// route publik: login & refresh
	authHandler := handlers.NewAuthHandler(db, ts, cfg)
	r.POST("/api/users/login", authHandler.Login)
	r.POST("/api/users/refresh", authHandler.Refresh)

	// protected dengan JWT
	jwtMW := middleware.NewJWTAuth(cfg)
	api := r.Group("/api", jwtMW)

	// logout (harus bawa AT valid), RT opsional
	api.POST("/users/logout", authHandler.Logout)

	// kategori
	catHandler := handlers.NewCategoryHandler(db)
	catGroup := api.Group("/categories")
	catHandler.Register(catGroup)

	// buku
	bookHandler := handlers.NewBookHandler(db)
	bookGroup := api.Group("/books")
	bookHandler.Register(bookGroup)

	return r
}
