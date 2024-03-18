package authDTOs

type LoginDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}