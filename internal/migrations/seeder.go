package migrations

import (
	"encoding/json"
	"os"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func Seeder(db *gorm.DB, logger *logrus.Logger) error {
	logger.Info("Seeding database...")

	return nil
}

func seedFromJSON[T any](filePath string, out *[]T, db *gorm.DB, log *logrus.Logger) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Warnf("Seed file not found: %s", filePath)
		return
	}

	if err := json.Unmarshal(data, out); err != nil {
		log.Warnf("Failed to parse JSON for %s: %v", filePath, err)
		return
	}

	var count int64
	if err := db.Model(out).Count(&count).Error; err != nil {
		log.Warnf("Failed to count records for %s: %v", filePath, err)
		return
	}

	if count == 0 {
		if err := db.Create(out).Error; err != nil {
			log.Warnf("Insert failed for %s: %v", filePath, err)
		} else {
			log.Infof("Inserted seed data from %s", filePath)
		}
	} else {
		log.Infof("Skipping insert for %s: table not empty", filePath)
	}
}
