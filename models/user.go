package models

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique_index; not null"`
	Password string `gorm:"not null"`
	Name     string `gorm:"not null"`
	Avatar   string
	Role     string `gorm:"default:'Member';not null"`
}

/*
//Method Hook

	func (u *User) BeforeSave(scope *gorm.Scope) {
		if u.Password == "" {
			return
		}

		hash, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
		scope.SetColumn("password", string(hash))
	}

	func (u *User) GenerateEncryptedPassword() string {
		hash, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
		return string(hash)
	}
*/
func (u *User) GenerateEncryptedPassword(pwd string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(pwd), 14)
	return string(hash)
}

func (u *User) Promote() {
	u.Role = "Editor"
}

func (u *User) Demote() {
	u.Role = "Member"
}
