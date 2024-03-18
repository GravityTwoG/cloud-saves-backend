package authDTOs

type ResetPasswordDTO struct {
	Token string `json:"token"`
	NewPassword string `json:"newPassword"`
}