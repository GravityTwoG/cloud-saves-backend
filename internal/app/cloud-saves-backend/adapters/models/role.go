package models

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/domain/user"
	"time"
)

type Role struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      user.RoleName `gorm:"unique;not null"`
}

func RoleFromEntity(role *user.Role) *Role {
	return &Role{
		ID:   role.GetId(),
		Name: role.GetName(),
	}
}

func RoleFromModel(roleModel *Role) *user.Role {
	return user.RoleFromDB(
		roleModel.ID,
		roleModel.Name,
	)
}
