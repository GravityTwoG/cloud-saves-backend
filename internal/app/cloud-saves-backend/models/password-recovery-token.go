package models

import (
	"gorm.io/gorm"
)

type PasswordRecoveryToken struct {
	gorm.Model
	Token string `gorm:"unique;not null"`
	UserID uint `gorm:"unique;not null"`
	User User `gorm:"onDelete:CASCADE"`
}