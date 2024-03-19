package authDTOs

type RequestPasswordResetDTO struct {
	Email string `json:"email" binding:"required,email"`
}
