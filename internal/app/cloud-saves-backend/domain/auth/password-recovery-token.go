package auth

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/domain/user"
	"time"

	"github.com/google/uuid"
)

type PasswordRecoveryToken struct {
	id        uint
	token     string
	user      user.User
	
	createdAt time.Time
	updatedAt time.Time
}

func NewPasswordRecoveryToken(user *user.User) *PasswordRecoveryToken {
	token := PasswordRecoveryToken{
		token: generateToken(),
		user:  *user,
	}
	return &token
}

func generateToken() string {
	return uuid.New().String()
}

func PasswordRecoveryTokenFromDB(
	ID uint,
	token string,
	user *user.User,
	createdAt time.Time,
	updatedAt time.Time,	
) *PasswordRecoveryToken {
	return &PasswordRecoveryToken{
		id:        ID,
		token:     token,
		user:      *user,
		
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (t *PasswordRecoveryToken) GetId() uint {
	return t.id
}

func (t *PasswordRecoveryToken) GetToken() string {
	return t.token
}

func (t *PasswordRecoveryToken) GetUser() *user.User {
	return &t.user
}

func (t *PasswordRecoveryToken) GetCreatedAt() time.Time {
	return t.createdAt
}

func (t *PasswordRecoveryToken) GetUpdatedAt() time.Time {
	return t.updatedAt
}

