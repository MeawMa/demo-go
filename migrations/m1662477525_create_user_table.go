package migrations

import (
	"demo-go/models"

	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

func m1662477525CreateUsersTable() *gormigrate.Migration {

	return &gormigrate.Migration{
		ID: "1662477525",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&models.User{}).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.DropTable("users").Error
		},
	}
}
