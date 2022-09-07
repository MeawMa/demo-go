package migrations

import (
	"demo-go/config"
	"log"

	"gopkg.in/gormigrate.v1"
)

func Migrate() {
	db := config.GetDB()
	m := gormigrate.New(
		db,
		gormigrate.DefaultOptions,
		[]*gormigrate.Migration{
			m1662403764CreateArticlesTable(),
			m1662462216CreateCatagoriesTable(),
			m1662469843AddCategoryToArticles(),
			m1662477525CreateUsersTable(),
			m1662485195AddUserIDToArticle(),
		},
	)

	if err := m.Migrate(); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}
}
