package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"backend/internal/config"
	"backend/internal/platform/database"
	"backend/modules/admin/domain"
)

type mongoOID struct {
	OID string `json:"$oid"`
}

type mongoDate struct {
	Date string `json:"$date"`
}

type mongoContragentType struct {
	ID        mongoOID  `json:"_id"`
	Name      string    `json:"name"`
	Icon      string    `json:"icon"`
	Status    string    `json:"status"`
	CreatedAt mongoDate `json:"createdAt"`
	UpdatedAt mongoDate `json:"updatedAt"`
}

func parseTime(raw string) time.Time {
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return time.Now()
	}
	return t
}

func main() {
	cfg := config.Load()
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("DB ulanishda xatolik: %v", err)
	}

	if err = db.AutoMigrate(&domain.ContragentType{}); err != nil {
		log.Fatalf("Migratsiyada xatolik: %v", err)
	}

	inputPath := "scripts/ttsa.contragenttypes.json"
	if len(os.Args) > 1 {
		inputPath = os.Args[1]
	}

	raw, err := os.ReadFile(inputPath)
	if err != nil {
		log.Fatalf("JSON faylni o'qishda xatolik: %v", err)
	}

	var items []mongoContragentType
	if err = json.Unmarshal(raw, &items); err != nil {
		log.Fatalf("JSON parse xatolik: %v", err)
	}

	createdCount := 0
	for _, item := range items {
		if item.Status == "" {
			item.Status = domain.StatusActive
		}

		var cnt int64
		if err = db.Model(&domain.ContragentType{}).Where("external_id = ?", item.ID.OID).Count(&cnt).Error; err != nil {
			log.Printf("Tekshirishda xatolik (%s): %v", item.Name, err)
			continue
		}
		if cnt > 0 {
			continue
		}

		row := domain.ContragentType{
			ExternalID: item.ID.OID,
			Name:       item.Name,
			Icon:       item.Icon,
			Status:     item.Status,
			CreatedAt:  parseTime(item.CreatedAt.Date),
			UpdatedAt:  parseTime(item.UpdatedAt.Date),
		}

		if err = db.Create(&row).Error; err != nil {
			log.Printf("Yozishda xatolik (%s): %v", item.Name, err)
			continue
		}
		createdCount++
	}

	log.Printf("Import yakunlandi. Yangi kontragent turlari: %d", createdCount)
}
