package controllers

import (
	"demo-go/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
)

type Categories struct {
	DB *gorm.DB
}
type categoryResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	Articles []struct {
		ID    uint   `json:"id"`
		Title string `json:"title"`
	} `json:"articles"`
}
type allCategoriesResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}
type createCategoryForm struct {
	Name string `form:"name" binding:"required"`
	Desc string `form:"desc" binding:"required"`
}
type updateCategoryForm struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

func (c *Categories) FindAll(ctx *gin.Context) {
	var categories []models.Category
	c.DB.Preload("Articles").Order("id desc").Find(&categories)

	var serializedCategories []allCategoriesResponse
	copier.Copy(&serializedCategories, &categories)

	ctx.JSON(http.StatusOK, gin.H{"categories": serializedCategories})
}

func (c *Categories) FindOne(ctx *gin.Context) {
	category, err := c.FindCategoryById(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var serializedCategory categoryResponse
	copier.Copy(&serializedCategory, &category)

	ctx.JSON(http.StatusOK, gin.H{"category": serializedCategory})
}
func (c *Categories) FindCategoryById(ctx *gin.Context) (*models.Category, error) {
	var category models.Category
	id := ctx.Param("id")

	if err := c.DB.Preload("Articles").First(&category, id).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

func (c *Categories) Create(ctx *gin.Context) {
	var form createCategoryForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var newCategory models.Category
	copier.Copy(&newCategory, &form)

	if err := c.DB.Create(&newCategory).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	var serializedCategory allCategoriesResponse
	copier.Copy(&serializedCategory, &newCategory)

	ctx.JSON(http.StatusCreated, gin.H{"category": serializedCategory})
}

func (c *Categories) Update(ctx *gin.Context) {
	var form updateCategoryForm

	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	category, err := c.FindCategoryById(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := c.DB.Model(&category).Update(&form).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var serializedCategory allCategoriesResponse
	copier.Copy(&serializedCategory, &category)

	ctx.JSON(http.StatusOK, gin.H{"category": serializedCategory})

}

func (c *Categories) Delete(ctx *gin.Context) {
	category, err := c.FindCategoryById(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := c.DB.Unscoped().Delete(&category).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}
