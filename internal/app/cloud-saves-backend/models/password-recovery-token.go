package models

import (
	"time"

	"github.com/google/uuid"
)

type PasswordRecoveryToken struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Token     string `gorm:"unique;not null"`
	UserID    uint   `gorm:"unique;not null"`
	User      User   `gorm:"onDelete:CASCADE"`
}

func NewPasswordRecoveryToken(user *User) *PasswordRecoveryToken {
	token := PasswordRecoveryToken{
		Token: generateToken(),
		User:  *user,
	}
	return &token
}

func generateToken() string {
	return uuid.New().String()
}
