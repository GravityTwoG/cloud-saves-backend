package user

import "cloud-saves-backend/internal/app/cloud-saves-backend/models"

type UserResponseDTO struct {
	Id        uint            `json:"id"`
	Email     string          `json:"email"`
	Username  string          `json:"username"`
	Role      models.RoleName `json:"role"`
	IsBlocked bool            `json:"isBlocked"`
}

func FromUser(user *models.User) *UserResponseDTO {
	return &UserResponseDTO{
		Id:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		Role:      user.Role.Name,
		IsBlocked: user.IsBlocked,
	}
}
