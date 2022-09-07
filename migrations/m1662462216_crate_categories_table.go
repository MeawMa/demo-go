package migrations

import (
	"demo-go/models"

	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

func m1662462216CreateCatagoriesTable() *gormigrate.Migration {

	return &gormigrate.Migration{
		ID: "1662462216",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&models.Category{}).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.DropTable("catagories").Error
		},
	}
}
