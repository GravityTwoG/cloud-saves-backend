package user

import "cloud-saves-backend/internal/app/cloud-saves-backend/models"

type UserResponseDTO struct {
	Id        uint            `json:"id"`
	Email     string          `json:"email"`
	Username  string          `json:"username"`
	Role      models.RoleName `json:"role"`
	IsBlocked bool            `json:"isBlocked"`
}
