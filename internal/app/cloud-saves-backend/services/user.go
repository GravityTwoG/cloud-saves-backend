package services

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/dto/common"
	"cloud-saves-backend/internal/app/cloud-saves-backend/dto/user"
	"cloud-saves-backend/internal/app/cloud-saves-backend/models"
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error

	Save(ctx context.Context, user *models.User) error

	GetByEmail(ctx context.Context, email string) (*models.User, error)

	GetByUsername(ctx context.Context, username string) (*models.User, error)

	GetById(ctx context.Context, userId uint) (*models.User, error)

	GetAll(
		ctx context.Context,
		dto common.GetResourceDTO,
	) (*common.ResourceDTO[models.User], error)
}

type UserService interface {
	GetUsers(
		ctx context.Context,
		dto common.GetResourceDTO,
	) (*user.UsersResponseDTO, error)
}

type userService struct {
	userRepo UserRepository
}

func NewUser(userRepo UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetUsers(
	ctx context.Context,
	dto common.GetResourceDTO,
) (*user.UsersResponseDTO, error) {
	users, err := s.userRepo.GetAll(ctx, dto)
	if err != nil {
		return nil, err
	}

	return user.FromUsers(users), nil
}
