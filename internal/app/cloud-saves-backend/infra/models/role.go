package models

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/domain/user"
	"time"
)


type RoleModel struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      user.RoleName `gorm:"unique;not null"`
}

func RoleFromEntity(role *user.Role) *RoleModel {
	return &RoleModel{
		ID:        role.GetId(),
		Name:      role.GetName(),
	}
}

func RoleFromModel(roleModel *RoleModel) *user.Role {
	return user.RoleFromDB(
		roleModel.ID, 
		roleModel.Name,
	)
}