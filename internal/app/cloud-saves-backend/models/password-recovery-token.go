package models

import "time"

type PasswordRecoveryToken struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Token     string `gorm:"unique;not null"`
	UserID    uint   `gorm:"unique;not null"`
	User      User   `gorm:"onDelete:CASCADE"`
}
