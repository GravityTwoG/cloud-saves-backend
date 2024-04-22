package user

import "cloud-saves-backend/internal/app/cloud-saves-backend/domain/user"

type UserDTO struct {
	Id        uint          `json:"id"`
	Email     string        `json:"email"`
	Username  string        `json:"username"`
	Role      user.RoleName `json:"role"`
	IsBlocked bool          `json:"isBlocked"`
}

func FromUser(user *user.User) *UserDTO {
	return &UserDTO{
		Id:        user.GetId(),
		Email:     user.GetEmail(),
		Username:  user.GetUsername(),
		Role:      user.GetRole().GetName(),
		IsBlocked: user.IsBlocked(),
	}
}
