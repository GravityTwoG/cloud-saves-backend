package user

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/domain/common"
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error

	Save(ctx context.Context, user *User) error

	GetByEmail(ctx context.Context, email string) (*User, error)

	GetByUsername(ctx context.Context, username string) (*User, error)

	GetById(ctx context.Context, userId uint) (*User, error)

	GetAll(
		ctx context.Context,
		dto common.GetResourceDTO,
	) (*common.ResourceDTO[User], error)
}