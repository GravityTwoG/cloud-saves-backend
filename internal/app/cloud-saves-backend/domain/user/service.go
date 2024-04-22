package user

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/domain/common"
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
	roleRepo RoleRepository
}

func NewUserService(userRepo UserRepository, roleRepo RoleRepository) UserService {
	return &userService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

func (s *userService) GetUsers(
	ctx context.Context,
	dto common.GetResourceDTO,
) (*common.ResourceDTO[User], error) {
	roleUser, err := s.roleRepo.GetByName(ctx, RoleUser)  
	if err != nil {
		return nil, err
	}
	users, err := s.userRepo.GetUsersWithRole(ctx, dto, roleUser)
	if err != nil {
		return nil, err
	}

	return users, nil
}
