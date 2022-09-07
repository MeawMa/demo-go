package controllers

import (
	"demo-go/config"
	"demo-go/models"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
)

type Users struct {
	DB *gorm.DB
}
type createUserForm struct {
	Email    string `json:"email"  binding:"required,email"`
	Password string `json:"password"  binding:"required,min=4"`
	Name     string `json:"name"  binding:"required"`
}
type updateUserForm struct {
	Email    string `json:"email"  binding:"omitempty,email"`
	Password string `json:"password"  binding:"omitempty,min=4"`
	Name     string `json:"name"`
}
type userResponse struct {
	ID     uint   `json:"id"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
	Name   string `json:"name"`
	Role   string `json:"role"`
}
type userPaging struct {
	Items  []userResponse `json:"itmes"`
	Paging *pagingResult  `json:"paging"`
}

func (u *Users) FindAll(ctx *gin.Context) {
	/*
		sub, _ := ctx.Get("sub")
		if sub.(*models.User).Role != "Admin" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Access Denied"})
			return
		}
	*/
	//var users []models.User
	users := []models.User{}

	query := u.DB.Order("id desc").Find(&users)
	term := ctx.Query("term")
	if term != "" {
		query = query.Where("name ILIKE ?", "%"+term+"%")
	}

	pagination := pagination{ctx: ctx, query: query, records: &users}
	paging := pagination.paginate()

	//var serializedUsers []userResponse
	serializedUsers := []userResponse{}
	copier.Copy(&serializedUsers, &users)

	ctx.JSON(http.StatusOK, gin.H{"user": userPaging{Items: serializedUsers, Paging: paging}})

}
func (u *Users) FindUserByID(ctx *gin.Context) (*models.User, error) {
	id := ctx.Param("id")
	var user models.User
	if err := u.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func (u *Users) FindOne(ctx *gin.Context) {
	user, err := u.FindUserByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusOK, gin.H{"user": serializedUser})
}
func (u *Users) Create(ctx *gin.Context) {
	var form createUserForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	var user models.User
	copier.Copy(&user, &form)
	user.Password = user.GenerateEncryptedPassword(user.Password)

	if err := u.DB.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusCreated, gin.H{"user": serializedUser})
}
func (u *Users) Update(ctx *gin.Context) {
	var form updateUserForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	user, err := u.FindUserByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if form.Password != "" {
		form.Password = user.GenerateEncryptedPassword(form.Password)
	}
	if err := u.DB.Model(&user).Update(&form).Error; err != nil {
		ctx.JSON(http.StatusUnavailableForLegalReasons, gin.H{"error": err.Error()})
		return
	}
	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusOK, gin.H{"user": serializedUser})
}
func (u *Users) Delete(ctx *gin.Context) {
	user, err := u.FindUserByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := u.DB.Unscoped().Delete(&user).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{})
}
func (u *Users) Promote(ctx *gin.Context) {
	user, err := u.FindUserByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	user.Promote()
	u.DB.Save(user)

	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusOK, gin.H{"user": serializedUser})
}
func (u *Users) Demote(ctx *gin.Context) {
	user, err := u.FindUserByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	user.Demote()
	u.DB.Save(user)

	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusOK, gin.H{"user": serializedUser})
}
func setUserImage(ctx *gin.Context, user *models.User) error {
	file, err := ctx.FormFile("avatar")
	if err != nil || file == nil {
		return err
	}

	if user.Avatar != "" {
		user.Avatar = strings.Replace(user.Avatar, os.Getenv("HOST"), "", 1)
		pwd, _ := os.Getwd()
		os.Remove(pwd + user.Avatar)
	}
	path := "uploads/users/" + strconv.Itoa(int(user.ID))
	os.MkdirAll(path, 0755)
	filename := path + "/" + file.Filename
	if err := ctx.SaveUploadedFile(file, filename); err != nil {
		return err
	}
	user.Avatar = os.Getenv("HOST") + "/" + filename

	db := config.GetDB()
	db.Save(user)

	return nil
}
