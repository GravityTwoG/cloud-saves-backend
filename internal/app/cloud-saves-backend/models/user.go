package models

import "time"

type User struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string `gorm:"unique;not null"`
	Email     string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	IsBlocked bool   `gorm:"not null"`
	RoleID    uint   `gorm:"not null"`
	Role      Role   `gorm:"onDelete:RESTRICT"`
}
