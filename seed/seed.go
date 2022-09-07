package seed

import (
	"demo-go/config"
	"demo-go/migrations"
	"demo-go/models"
	"math/rand"
	"strconv"

	"github.com/bxcodec/faker/v4"
	"github.com/labstack/gommon/log"
)

func Load() {
	db := config.GetDB()
	db.DropTableIfExists("articles", "categories", "migrations", "users")
	migrations.Migrate()

	log.Info("Creating admin...")
	admin := models.User{
		Email:    "admin@admin.com",
		Password: "adminadmin",
		Name:     "Admin Admin",
		Role:     "Admin",
	}
	admin.Password = admin.GenerateEncryptedPassword(admin.Password)
	db.Create(&admin)

	log.Info("Craeting users...")
	numOfUsers := 50
	users := make([]models.User, 0, numOfUsers)
	userRoles := [2]string{"Member", "Editor"}
	for u := 1; u <= numOfUsers; u++ {
		user := models.User{
			Email:    faker.Email(),
			Password: "useruser",
			Name:     faker.Name(),
			Role:     userRoles[rand.Intn(2)],
		}
		user.Password = user.GenerateEncryptedPassword(user.Password)
		db.Create(&user)
		users = append(users, user)
	}

	log.Info("Creating categories....")
	numOfCategories := 20
	categories := make([]models.Category, 0, numOfCategories)
	for i := 1; i <= numOfCategories; i++ {
		category := models.Category{
			Name: faker.Word(),
			Desc: faker.Paragraph(),
		}
		db.Create(&category)
		categories = append(categories, category)
	}

	log.Info("Creating articles....")

	numOfArticles := 50
	articles := make([]models.Article, 0, numOfArticles)
	for a := 1; a <= numOfArticles; a++ {
		article := models.Article{
			Title:      faker.Sentence(),
			Excerpt:    faker.Sentence(),
			Image:      "https://source.unsplash.com/random/100x100?" + strconv.Itoa(a),
			CategoryID: uint(rand.Intn(numOfCategories)) + 1,
			UserID:     uint(rand.Intn(numOfUsers) + 1),
		}

		db.Create(&article)
		articles = append(articles, article)
	}

}
