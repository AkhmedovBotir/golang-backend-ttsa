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

type mongoLocation struct {
	ID        mongoOID  `json:"_id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Parent    *mongoOID `json:"parent"`
	Code      string    `json:"code"`
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

	if err = db.AutoMigrate(&domain.Region{}, &domain.District{}, &domain.MFY{}); err != nil {
		log.Fatalf("Migratsiyada xatolik: %v", err)
	}

	inputPath := "scripts/ttsa.regions.json"
	if len(os.Args) > 1 {
		inputPath = os.Args[1]
	}

	raw, err := os.ReadFile(inputPath)
	if err != nil {
		log.Fatalf("JSON faylni o'qishda xatolik: %v", err)
	}

	var items []mongoLocation
	if err = json.Unmarshal(raw, &items); err != nil {
		log.Fatalf("JSON parse xatolik: %v", err)
	}

	regionMap := map[string]uint{}
	districtMap := map[string]uint{}

	var regionCount, districtCount, mfyCount int

	for _, item := range items {
		if item.Type != "region" {
			continue
		}

		row := domain.Region{
			ExternalID: item.ID.OID,
			Name:       item.Name,
			Code:       item.Code,
			Status:     item.Status,
			CreatedAt:  parseTime(item.CreatedAt.Date),
			UpdatedAt:  parseTime(item.UpdatedAt.Date),
		}

		if row.Status == "" {
			row.Status = domain.StatusActive
		}

		var existingID uint
		var existingCount int64
		if err = db.Model(&domain.Region{}).Where("external_id = ?", row.ExternalID).Count(&existingCount).Error; err != nil {
			log.Printf("Viloyat tekshirishda xatolik (%s): %v", item.Name, err)
			continue
		}
		if existingCount > 0 {
			if err = db.Model(&domain.Region{}).Where("external_id = ?", row.ExternalID).Select("id").Scan(&existingID).Error; err == nil {
				regionMap[item.ID.OID] = existingID
			}
			continue
		}

		if err = db.Create(&row).Error; err != nil {
			log.Printf("Viloyat yozishda xatolik (%s): %v", item.Name, err)
			continue
		}
		regionMap[item.ID.OID] = row.ID
		regionCount++
	}

	for _, item := range items {
		if item.Type != "district" {
			continue
		}
		if item.Parent == nil {
			continue
		}
		regionID, ok := regionMap[item.Parent.OID]
		if !ok {
			continue
		}

		row := domain.District{
			ExternalID: item.ID.OID,
			RegionID:   regionID,
			Name:       item.Name,
			Code:       item.Code,
			Status:     item.Status,
			CreatedAt:  parseTime(item.CreatedAt.Date),
			UpdatedAt:  parseTime(item.UpdatedAt.Date),
		}
		if row.Status == "" {
			row.Status = domain.StatusActive
		}

		var existingID uint
		var existingCount int64
		if err = db.Model(&domain.District{}).Where("external_id = ?", row.ExternalID).Count(&existingCount).Error; err != nil {
			log.Printf("Tuman tekshirishda xatolik (%s): %v", item.Name, err)
			continue
		}
		if existingCount > 0 {
			if err = db.Model(&domain.District{}).Where("external_id = ?", row.ExternalID).Select("id").Scan(&existingID).Error; err == nil {
				districtMap[item.ID.OID] = existingID
			}
			continue
		}

		if err = db.Create(&row).Error; err != nil {
			log.Printf("Tuman yozishda xatolik (%s): %v", item.Name, err)
			continue
		}
		districtMap[item.ID.OID] = row.ID
		districtCount++
	}

	for _, item := range items {
		if item.Type != "mfy" {
			continue
		}
		if item.Parent == nil {
			continue
		}
		districtID, ok := districtMap[item.Parent.OID]
		if !ok {
			continue
		}

		row := domain.MFY{
			ExternalID: item.ID.OID,
			DistrictID: districtID,
			Name:       item.Name,
			Code:       item.Code,
			Status:     item.Status,
			CreatedAt:  parseTime(item.CreatedAt.Date),
			UpdatedAt:  parseTime(item.UpdatedAt.Date),
		}
		if row.Status == "" {
			row.Status = domain.StatusActive
		}

		var existingCount int64
		if err = db.Model(&domain.MFY{}).Where("external_id = ?", row.ExternalID).Count(&existingCount).Error; err != nil {
			log.Printf("MFY tekshirishda xatolik (%s): %v", item.Name, err)
			continue
		}
		if existingCount > 0 {
			continue
		}

		if err = db.Create(&row).Error; err != nil {
			log.Printf("MFY yozishda xatolik (%s): %v", item.Name, err)
			continue
		}
		mfyCount++
	}

	log.Printf("Import yakunlandi. Yangi yozuvlar: region=%d, district=%d, mfy=%d", regionCount, districtCount, mfyCount)
}
