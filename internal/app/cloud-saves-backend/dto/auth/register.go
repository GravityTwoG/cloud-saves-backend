package auth

type RegisterDTO struct {
	Email    string `json:"email" binding:"required,email,max=256"`
	Username string `json:"username" binding:"required,alphanum,min=3,max=32"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}
