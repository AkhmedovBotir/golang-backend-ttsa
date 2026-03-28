package main

import (
	"log"

	"backend/internal/config"
	"backend/internal/pkg/security"
	"backend/internal/platform/database"
	"backend/modules/admin/domain"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("DB ulanishda xatolik: %v", err)
	}

	if err = db.AutoMigrate(&domain.Admin{}); err != nil {
		log.Fatalf("Migratsiyada xatolik: %v", err)
	}

	var existing domain.Admin
	err = db.Where("username = ? OR phone = ?", cfg.GeneralAdminUsername, cfg.GeneralAdminPhone).First(&existing).Error
	if err == nil {
		log.Println("General admin allaqachon mavjud.")
		return
	}

	hashedPassword, err := security.HashPassword(cfg.GeneralAdminPassword)
	if err != nil {
		log.Fatalf("Parolni hash qilishda xatolik: %v", err)
	}

	newAdmin := domain.Admin{
		Name:     cfg.GeneralAdminName,
		Role:     cfg.GeneralAdminRole,
		Phone:    cfg.GeneralAdminPhone,
		Username: cfg.GeneralAdminUsername,
		Password: hashedPassword,
		Status:   cfg.GeneralAdminStatus,
	}

	if err = db.Create(&newAdmin).Error; err != nil {
		log.Fatalf("General admin yaratishda xatolik: %v", err)
	}

	log.Println("General admin muvaffaqiyatli yaratildi.")
}
