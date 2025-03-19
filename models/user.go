package models

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User struct with validation tags
type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()" json:"id"`
	Name     string    `gorm:"size:100;not null;" json:"name" validate:"required,min=4" json:"name"`
	Username string    `gorm:"size:100;uniqueIndex;not null" json:"username" validate:"required,min=4" json:"username"`
	Email    string    `gorm:"size:100;not null;unique" json:"email" validate:"required,email" json:"email"`
	Password string    `gorm:"size:100;not null" json:"password" validate:"required,min=6" json:"password"`
	Avatar   string    `gorm:"type:text" json:"avatar"`
}

// Validate function to check struct fields
func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}
