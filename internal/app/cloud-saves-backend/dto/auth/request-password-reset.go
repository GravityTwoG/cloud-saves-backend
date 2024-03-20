package auth

type RequestPasswordResetDTO struct {
	Email string `json:"email" binding:"required,email"`
}
