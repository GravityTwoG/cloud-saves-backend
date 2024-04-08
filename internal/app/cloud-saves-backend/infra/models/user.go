package models

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/domain/user"
	"time"
)

type UserModel struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string `gorm:"unique;not null"`
	Email     string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	IsBlocked bool   `gorm:"not null"`
	RoleID    uint   `gorm:"not null"`
	Role      RoleModel   `gorm:"onDelete:RESTRICT"`
}

func UserFromEntity(user *user.User) *UserModel {
	return &UserModel{
		ID:        user.GetId(),
		Username:  user.GetUsername(),
		Email:     user.GetEmail(),
		Password:  user.GetPassword(),
		IsBlocked: user.IsBlocked(),
		RoleID:    user.GetRole().GetId(),
		Role:      *RoleFromEntity(user.GetRole()),
	}
}

func UserFromModel(userModel *UserModel) *user.User {
	return user.UserFromDB(
		userModel.ID, 
		userModel.Username, 
		userModel.Email, 
		userModel.Password, 
		RoleFromModel(&userModel.Role), 
		userModel.IsBlocked,
	)
}
