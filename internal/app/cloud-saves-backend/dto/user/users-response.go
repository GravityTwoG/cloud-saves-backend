package user

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/dto/common"
	"cloud-saves-backend/internal/app/cloud-saves-backend/models"
)

type UsersResponseDTO struct {
	Items      []UserResponseDTO `json:"items"`
	TotalCount int               `json:"totalCount"`
}

func FromUsers(users *common.ResourceDTO[models.User]) *UsersResponseDTO {
	usersDto := &UsersResponseDTO{
		Items:      make([]UserResponseDTO, len(users.Items)),
		TotalCount: users.TotalCount,
	}

	for i, userModel := range users.Items {
		usersDto.Items[i] = *FromUser(&userModel)
	}

	return usersDto
}
