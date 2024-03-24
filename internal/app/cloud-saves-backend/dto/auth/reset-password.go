package auth

type ResetPasswordDTO struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}
