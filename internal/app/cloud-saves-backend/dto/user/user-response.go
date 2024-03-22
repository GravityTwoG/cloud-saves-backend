package user

import "cloud-saves-backend/internal/app/cloud-saves-backend/models"

type UserResponseDTO struct {
	Id        uint            `json:"id"`
	Email     string          `json:"email"`
	Username  string          `json:"username"`
	Role      models.RoleName `json:"role"`
	IsBlocked bool            `json:"isBlocked"`
}

func (u *UserResponseDTO) FromUser(user *models.User) *UserResponseDTO {
	u.Id = user.ID
	u.Email = user.Email
	u.Username = user.Username
	u.Role = user.Role.Name
	u.IsBlocked = user.IsBlocked
	return u
}
