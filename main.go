package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"demo-go/config"
	"demo-go/migrations"
	"demo-go/routes"

	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	config.InitDB()
	defer config.CloseDB()
	migrations.Migrate()
	//seed.Load()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("Authorization")

	router := gin.Default()
	router.Use(cors.New(corsConfig))
	router.Static("/uploads", "./uploads")
	uploadDirs := [...]string{"articles", "users"}
	for _, dir := range uploadDirs {
		os.MkdirAll("uploads/"+dir, 0755)
	}

	routes.Serve(router)

	router.Run(":" + os.Getenv("PORT"))
}
