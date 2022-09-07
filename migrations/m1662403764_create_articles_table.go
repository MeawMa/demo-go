package migrations

import (
	"demo-go/models"

	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

func m1662403764CreateArticlesTable() *gormigrate.Migration {

	return &gormigrate.Migration{
		ID: "1662403764",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&models.Article{}).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.DropTable("articles").Error
		},
	}
}
