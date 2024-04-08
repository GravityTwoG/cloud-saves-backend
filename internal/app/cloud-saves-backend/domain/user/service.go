package user

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/dto/common"
	"context"
)

type UserService interface {
	GetUsers(
		ctx context.Context,
		dto common.GetResourceDTO,
	) (*common.ResourceDTO[User], error)
}

type userService struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetUsers(
	ctx context.Context,
	dto common.GetResourceDTO,
) (*common.ResourceDTO[User], error) {
	users, err := s.userRepo.GetAll(ctx, dto)
	if err != nil {
		return nil, err
	}

	return users, nil
}
