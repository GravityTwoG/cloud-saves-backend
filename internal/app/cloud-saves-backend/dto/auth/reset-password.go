package authDTOs

type ResetPasswordDTO struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=8,max=64"`
}