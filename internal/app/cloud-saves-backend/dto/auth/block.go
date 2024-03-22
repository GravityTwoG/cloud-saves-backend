package auth

type BlockUserDTO struct {
	UserId uint `uri:"userId" binding:"required"`
}
