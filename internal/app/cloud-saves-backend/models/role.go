package models

import "time"

type RoleName string

const RoleUser RoleName = "ROLE_USER"
const RoleAdmin RoleName = "ROLE_ADMIN"

type Role struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      RoleName `gorm:"unique;not null"`
}
