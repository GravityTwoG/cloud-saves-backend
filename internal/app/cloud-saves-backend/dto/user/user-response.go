package userDTOs

type UserResponseDTO struct {
	Id uint `json:"id"`
	Email string `json:"email"`
	Username string `json:"username"`
	Role string `json:"role"`
	IsBlocked bool `json:"isBlocked"`
}