package main

import (
	"log"

	"backend/internal/config"
	"backend/internal/platform/database"
	"backend/internal/platform/httpserver"
	"backend/modules/admin"
	"backend/modules/contragents"
	"backend/modules/marketplace"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	log.Println("Server ishga tushdi")

	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("DB ulanishda xatolik: %v", err)
	}

	router := httpserver.New()
	if err = admin.RegisterRoutes(router, db, cfg.JWTSecret, cfg.JWTExpireHours); err != nil {
		log.Fatalf("Router ulashda xatolik: %v", err)
	}
	if err = contragents.RegisterRoutes(router, db, cfg.JWTSecret, cfg.JWTExpireHours); err != nil {
		log.Fatalf("Contragent auth router ulashda xatolik: %v", err)
	}
	if err = marketplace.RegisterRoutes(router, db, cfg.JWTSecret, cfg.JWTExpireHours); err != nil {
		log.Fatalf("Marketplace auth router ulashda xatolik: %v", err)
	}

	if err = router.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("Server ishga tushmadi: %v", err)
	}
}
