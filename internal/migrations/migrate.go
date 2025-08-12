package migrations

import (
	"golectro-payment/internal/entity"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(entity.Invoice{})
}
