package auth

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/domain/user"
	"context"
)

type RoleRepository interface {
	GetByName(ctx context.Context, name user.RoleName) (*user.Role, error)
}

type PasswordRecoveryTokenRepository interface {
	Create(ctx context.Context, token *PasswordRecoveryToken) error

	Save(ctx context.Context, token *PasswordRecoveryToken) error

	GetByToken(ctx context.Context, token string) (*PasswordRecoveryToken, error)

	GetByUserId(ctx context.Context, userId uint) (*PasswordRecoveryToken, error)

	Delete(ctx context.Context, token *PasswordRecoveryToken) error
}
