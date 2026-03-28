package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"backend/internal/config"
	"backend/internal/platform/database"
	"backend/modules/admin/domain"
	"gorm.io/gorm"
)

type mongoOID struct {
	OID string `json:"$oid"`
}

type mongoDate struct {
	Date string `json:"$date"`
}

type mongoCategory struct {
	ID        mongoOID  `json:"_id"`
	Name      string    `json:"name"`
	Image     *string   `json:"image"`
	Censored  bool      `json:"censored"`
	Parent    *mongoOID `json:"parent"`
	Status    string    `json:"status"`
	CreatedAt mongoDate `json:"createdAt"`
	UpdatedAt mongoDate `json:"updatedAt"`
	Slug      string    `json:"slug"`
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

	if err = db.AutoMigrate(&domain.Category{}); err != nil {
		log.Fatalf("Migratsiyada xatolik: %v", err)
	}

	inputPath := "scripts/ttsa.categories.json"
	if len(os.Args) > 1 {
		inputPath = os.Args[1]
	}

	raw, err := os.ReadFile(inputPath)
	if err != nil {
		log.Fatalf("JSON faylni o'qishda xatolik: %v", err)
	}

	var items []mongoCategory
	if err = json.Unmarshal(raw, &items); err != nil {
		log.Fatalf("JSON parse xatolik: %v", err)
	}

	idMap := map[string]uint{}
	createdCount := 0
	updatedCount := 0

	// 1-pass: parent_id bo'lmaganlar
	for _, item := range items {
		if item.Parent != nil {
			continue
		}
		upsertCategory(db, item, nil, idMap, &createdCount, &updatedCount)
	}

	// 2-pass: parent_id borlar
	for _, item := range items {
		if item.Parent == nil {
			continue
		}
		parentID, ok := idMap[item.Parent.OID]
		if !ok {
			log.Printf("Parent topilmadi, skip qilindi: %s (%s)", item.Name, item.ID.OID)
			continue
		}
		upsertCategory(db, item, &parentID, idMap, &createdCount, &updatedCount)
	}

	log.Printf("Import yakunlandi. Created=%d, Updated=%d", createdCount, updatedCount)
}

func upsertCategory(
	db *gorm.DB,
	item mongoCategory,
	parentID *uint,
	idMap map[string]uint,
	createdCount *int,
	updatedCount *int,
) {
	status := item.Status
	if status == "" {
		status = domain.StatusActive
	}
	image := ""
	if item.Image != nil {
		image = *item.Image
	}

	var existing domain.Category
	err := db.Where("external_id = ?", item.ID.OID).First(&existing).Error
	if err == nil {
		existing.Name = item.Name
		existing.Slug = item.Slug
		existing.Image = image
		existing.Censored = item.Censored
		existing.ParentID = parentID
		existing.Status = status
		existing.UpdatedAt = parseTime(item.UpdatedAt.Date)
		if err = db.Save(&existing).Error; err != nil {
			log.Printf("Yangilashda xatolik (%s): %v", item.Name, err)
			return
		}
		idMap[item.ID.OID] = existing.ID
		*updatedCount = *updatedCount + 1
		return
	}
	if err != gorm.ErrRecordNotFound {
		log.Printf("Tekshirishda xatolik (%s): %v", item.Name, err)
		return
	}

	row := domain.Category{
		ExternalID: item.ID.OID,
		Name:       item.Name,
		Slug:       item.Slug,
		Image:      image,
		Censored:   item.Censored,
		ParentID:   parentID,
		Status:     status,
		CreatedAt:  parseTime(item.CreatedAt.Date),
		UpdatedAt:  parseTime(item.UpdatedAt.Date),
	}
	if err = db.Create(&row).Error; err != nil {
		log.Printf("Yozishda xatolik (%s): %v", item.Name, err)
		return
	}
	idMap[item.ID.OID] = row.ID
	*createdCount = *createdCount + 1
}
