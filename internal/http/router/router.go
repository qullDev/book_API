package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// endpoint health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	return r
}
