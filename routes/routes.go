package routes

import (
	"demo-go/config"
	"demo-go/controllers"
	"demo-go/middleware"

	"github.com/gin-gonic/gin"
)

func Serve(r *gin.Engine) {
	db := config.GetDB()

	authenticate := middleware.Authenticate().MiddlewareFunc()
	authorize := middleware.Authorize()

	V1 := r.Group("/api/v1")
	articlesGroup := V1.Group("articles")
	articleController := controllers.Articles{DB: db}
	r.GET("", articleController.Hello)
	articlesGroup.GET("", articleController.FindAll)
	articlesGroup.GET("/:id", articleController.FindOne)
	articlesGroup.Use(authenticate, authorize)
	{
		articlesGroup.POST("", authenticate, articleController.CreateArticle)
		articlesGroup.PATCH("/:id", articleController.UpdateArticle)
		articlesGroup.DELETE("/:id", articleController.Delete)
	}

	categoriesGroup := V1.Group("categories")
	categoryController := controllers.Categories{DB: db}
	categoriesGroup.GET("", categoryController.FindAll)
	categoriesGroup.GET("/:id", categoryController.FindOne)
	categoriesGroup.Use(authenticate, authorize)
	{
		categoriesGroup.POST("", categoryController.Create)
		categoriesGroup.PATCH("/:id", categoryController.Update)
		categoriesGroup.DELETE("/:id", categoryController.Delete)
	}
	authGrop := V1.Group("auth")
	authController := controllers.Auth{DB: db}
	{
		authGrop.POST("/sign-up", authController.Signup)
		authGrop.POST("/sign-in", middleware.Authenticate().LoginHandler)
		authGrop.GET("/profile", authenticate, authController.GetProfile)
		authGrop.PATCH("/profile", authenticate, authController.UpdateProfile)
	}
	userGrop := V1.Group("users")
	userController := controllers.Users{DB: db}
	userGrop.Use(authenticate, authorize)
	{
		userGrop.GET("", userController.FindAll)
		userGrop.GET("/:id", userController.FindOne)
		userGrop.POST("", userController.Create)
		userGrop.PATCH("/:id", userController.Update)
		userGrop.DELETE("/:id", userController.Delete)
		userGrop.PATCH("/:id/promote", userController.Promote)
		userGrop.PATCH("/:id/demote", userController.Demote)
	}
	dashboardController := controllers.Dashboard{DB: db}
	dashboardGroup := V1.Group("dashboard")
	dashboardGroup.Use(authenticate, authorize)
	{
		dashboardGroup.GET("", dashboardController.GetInfo)
	}

}
