package migrations

import (
	"demo-go/models"

	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

func m1662485195AddUserIDToArticle() *gormigrate.Migration {

	return &gormigrate.Migration{
		ID: "1662485195",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&models.Article{}).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Model(&models.Article{}).DropColumn("user_id").Error
		},
	}
}
