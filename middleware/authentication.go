package middleware

import (
	"demo-go/config"
	"demo-go/models"
	"log"
	"os"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type login struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=4"`
}

var identityKey = "sub"

func Authenticate() *jwt.GinJWTMiddleware {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Key:           []byte(os.Getenv("SECRET_KEY")),
		Timeout:       time.Hour,
		IdentityKey:   identityKey,
		TokenLookup:   "header: Authorization",
		TokenHeadName: "Bearer",
		IdentityHandler: func(ctx *gin.Context) interface{} {
			var user models.User
			claims := jwt.ExtractClaims(ctx)
			id := claims[identityKey]
			db := config.GetDB()
			if db.First(&user, uint(id.(float64))).RecordNotFound() {
				return nil
			}
			return &user
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var form login
			var user models.User
			if err := c.ShouldBindJSON(&form); err != nil {
				return nil, jwt.ErrMissingLoginValues
			}
			db := config.GetDB()
			if db.Where("email = ?", form.Email).First(&user).RecordNotFound() {
				return nil, jwt.ErrFailedAuthentication
			}
			if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)); err != nil {
				return nil, jwt.ErrFailedAuthentication
			}
			return &user, nil
		},
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.User); ok {
				claims := jwt.MapClaims{
					identityKey: v.ID,
				}
				return claims
			}
			return jwt.MapClaims{}
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{"error": message})
		},
	})
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}
	return authMiddleware
}
