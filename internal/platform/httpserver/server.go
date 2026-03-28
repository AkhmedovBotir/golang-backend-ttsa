package httpserver

import (
	"net/http"
	"time"

	"backend/internal/pkg/response"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:8082", "http://localhost:8083"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/health", func(c *gin.Context) {
		response.JSON(c, http.StatusOK, "Server ishlayapti", gin.H{"status": "ok"}, nil)
	})

	router.NoRoute(func(c *gin.Context) {
		response.JSON(c, http.StatusNotFound, "Bunday endpoint topilmadi", nil, nil)
	})

	router.NoMethod(func(c *gin.Context) {
		response.JSON(c, http.StatusMethodNotAllowed, "Bu HTTP metodiga ruxsat yo'q", nil, nil)
	})

	return router
}
