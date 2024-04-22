package auth

import (
	"context"
)



type PasswordRecoveryTokenRepository interface {
	Create(ctx context.Context, token *PasswordRecoveryToken) error

	Save(ctx context.Context, token *PasswordRecoveryToken) error

	GetByToken(ctx context.Context, token string) (*PasswordRecoveryToken, error)

	GetByUserId(ctx context.Context, userId uint) (*PasswordRecoveryToken, error)

	Delete(ctx context.Context, token *PasswordRecoveryToken) error
}
