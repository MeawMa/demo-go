package controllers

import (
	"demo-go/models"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
)

type Articles struct {
	DB *gorm.DB
}

type articleResponse struct {
	ID         uint   `json:"id"`
	Title      string `json:"title"`
	Excerpt    string `json:"excerpt"`
	Body       string `json:"body"`
	Image      string `json:"image"`
	CategoryID uint   `json:"categoryId"`
	Category   struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"category"`
	User struct {
		Name   string `json:"name"`
		Avatar string `json:"avatar"`
	} `json:"user"`
}
type createOrUpdateArticleResponse struct {
	ID         uint   `json:"id"`
	Title      string `json:"title"`
	Excerpt    string `json:"excerpt"`
	Body       string `json:"body"`
	Image      string `json:"image"`
	CategoryID uint   `json:"categoryId"`
	UserID     uint   `json:"userId"`
}

type createArticleForm struct {
	Title      string                `form:"title" binding:"required"`
	Body       string                `form:"body" binding:"required"`
	Excerpt    string                `form:"excerpt" binding:"required"`
	CategoryID uint                  `form:"categoryId" binding:"required"`
	Image      *multipart.FileHeader `form:"image"`
}
type updateArticleForm struct {
	Title   string                `form:"title"`
	Body    string                `form:"body"`
	Excerpt string                `form:"excerpt"`
	Image   *multipart.FileHeader `form:"image"`
}
type articlesPaging struct {
	Items  []articleResponse `json:"itmes"`
	Paging *pagingResult     `json:"paging"`
}

func (a *Articles) FindAll(ctx *gin.Context) {
	//var articles []models.Article nill slice null
	articles := []models.Article{} // empty slice []

	query := a.DB.Preload("User").Preload("Category").Order("id desc")
	categoryId := ctx.Query("categoryId")
	if categoryId != "" {
		query = query.Where("category_id = ?", categoryId)
	}
	term := ctx.Query("term")
	if term != "" {
		query = query.Where("title ILIKE ?", "%"+term+"%")
	}
	pagination := pagination{ctx: ctx, query: query, records: &articles}
	paging := pagination.paginate()

	//var serializedArticle []articleResponse // nill slice null
	serializedArticle := []articleResponse{} // empty slice []
	copier.Copy(&serializedArticle, &articles)

	ctx.JSON(http.StatusOK, gin.H{"articles": articlesPaging{Items: serializedArticle, Paging: paging}})
}
func (a *Articles) FindOne(ctx *gin.Context) {
	article, err := a.FindArticleById(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	serializedArticle := articleResponse{}
	copier.Copy(&serializedArticle, article)
	ctx.JSON(http.StatusOK, gin.H{"article": serializedArticle})

}
func (a *Articles) FindArticleById(ctx *gin.Context) (*models.Article, error) {
	var article models.Article
	id := ctx.Param("id")

	if err := a.DB.Preload("User").Preload("Category").First(&article, id).Error; err != nil {
		return nil, err
	}

	return &article, nil

}
func (a *Articles) CreateArticle(ctx *gin.Context) {
	var form createArticleForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	user, _ := ctx.Get("sub")

	var newArticle models.Article
	copier.Copy(&newArticle, &form)

	newArticle.User = *user.(*models.User)

	if err := a.DB.Create(&newArticle).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	a.setArticleImage(ctx, &newArticle)
	serializedArticle := createOrUpdateArticleResponse{}
	copier.Copy(&serializedArticle, newArticle)

	ctx.JSON(http.StatusCreated, gin.H{"article": serializedArticle})

}

func (a *Articles) UpdateArticle(ctx *gin.Context) {
	var form updateArticleForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	article, err := a.FindArticleById(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := a.DB.Model(&article).Update(&form).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	a.setArticleImage(ctx, article)

	var serializedArticle createOrUpdateArticleResponse
	copier.Copy(&serializedArticle, article)

	ctx.JSON(http.StatusOK, gin.H{"article": serializedArticle})

}

func (a *Articles) DeleteArticle(ctx *gin.Context) {
	article, err := a.FindArticleById(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := a.DB.Delete(&article).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (a *Articles) Delete(ctx *gin.Context) {
	article, err := a.FindArticleById(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := a.DB.Unscoped().Delete(&article).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (a *Articles) setArticleImage(ctx *gin.Context, article *models.Article) error {
	file, err := ctx.FormFile("image")
	if err != nil || file == nil {
		return err
	}

	if article.Image != "" {
		article.Image = strings.Replace(article.Image, os.Getenv("HOST"), "", 1)
		pwd, _ := os.Getwd()
		os.Remove(pwd + article.Image)
	}
	path := "uploads/articles/" + strconv.Itoa(int(article.ID))
	os.MkdirAll(path, 0755)
	filename := path + "/" + file.Filename
	if err := ctx.SaveUploadedFile(file, filename); err != nil {
		return err
	}
	article.Image = os.Getenv("HOST") + "/" + filename
	a.DB.Save(article)

	return nil
}

func (a *Articles) Hello(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Hello, GO!")
}
