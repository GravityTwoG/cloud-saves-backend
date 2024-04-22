package user

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/domain/common"
	"cloud-saves-backend/internal/app/cloud-saves-backend/domain/user"
)

type UsersResponseDTO struct {
	Items      []UserDTO `json:"items"`
	TotalCount int       `json:"totalCount"`
}

func FromUsers(users *common.ResourceDTO[user.User]) *UsersResponseDTO {
	usersDto := &UsersResponseDTO{
		Items:      make([]UserDTO, len(users.Items)),
		TotalCount: users.TotalCount,
	}

	for i, user := range users.Items {
		usersDto.Items[i] = *FromUser(&user)
	}

	return usersDto
}
