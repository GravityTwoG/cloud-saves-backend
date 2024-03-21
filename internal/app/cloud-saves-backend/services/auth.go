package services

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/dto/auth"
	userDTOs "cloud-saves-backend/internal/app/cloud-saves-backend/dto/user"
	"cloud-saves-backend/internal/app/cloud-saves-backend/models"
	password_utils "cloud-saves-backend/internal/app/cloud-saves-backend/utils/password-utils"
	"context"
	"fmt"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
)

type RoleRepository interface {
	GetByName(ctx context.Context, name string) (*models.Role, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error

	Save(ctx context.Context, user *models.User) error

	GetByEmail(ctx context.Context, email string) (*models.User, error)

	GetByUsername(ctx context.Context, username string) (*models.User, error)

	GetById(ctx context.Context, userId uint) (*models.User, error)
}

type PasswordRecoveryTokenRepository interface {
	Create(ctx context.Context, token *models.PasswordRecoveryToken) error

	Save(ctx context.Context, token *models.PasswordRecoveryToken) error

	GetByToken(ctx context.Context, token string) (*models.PasswordRecoveryToken, error)

	GetByUserId(ctx context.Context, userId uint) (*models.PasswordRecoveryToken, error)

	Delete(ctx context.Context, token *models.PasswordRecoveryToken) error
}

type AuthService interface {
	Register(dto *auth.RegisterDTO) (*userDTOs.UserResponseDTO, error)

	Login(dto *auth.LoginDTO) (*userDTOs.UserResponseDTO, error)

	ChangePassword(userId uint, dto *auth.ChangePasswordDTO) error

	RequestPasswordReset(dto *auth.RequestPasswordResetDTO) error

	ResetPassword(dto *auth.ResetPasswordDTO) error
}

type authService struct {
	trManager trm.Manager
	context   context.Context

	roleRepo     RoleRepository
	userRepo     UserRepository
	recoveryRepo PasswordRecoveryTokenRepository

	emailService EmailService
}

func NewAuth(trManager trm.Manager, context context.Context, roleRepo RoleRepository, userRepo UserRepository, recoveryRepo PasswordRecoveryTokenRepository, emailService EmailService) AuthService {
	return &authService{
		trManager:    trManager,
		context:      context,
		roleRepo:     roleRepo,
		userRepo:     userRepo,
		recoveryRepo: recoveryRepo,
		emailService: emailService,
	}
}

func (s *authService) Register(registerDTO *auth.RegisterDTO) (*userDTOs.UserResponseDTO, error) {
	roleUser, err := s.roleRepo.GetByName(s.context, "ROLE_USER")
	if err != nil {
		return nil, err
	}

	hashedPassword, err := password_utils.HashPassword(registerDTO.Password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Username:  registerDTO.Username,
		Email:     registerDTO.Email,
		Password:  hashedPassword,
		IsBlocked: false,
		Role:      *roleUser,
	}

	err = s.userRepo.Create(s.context, &user)
	if err != nil {
		return nil, err
	}

	userResponseDTO := userDTOs.UserResponseDTO{
		Id:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		Role:      user.Role.Name,
		IsBlocked: user.IsBlocked,
	}

	return &userResponseDTO, nil
}

func (s *authService) Login(
	loginDTO *auth.LoginDTO,
) (*userDTOs.UserResponseDTO, error) {
	user, err := s.userRepo.GetByUsername(s.context, loginDTO.Username)
	if err != nil {
		return nil, err
	}

	if !password_utils.ComparePasswords(user.Password, loginDTO.Password) {
		return nil, fmt.Errorf("INCORRECT_USERNAME_OR_PASSWORD")
	}
	if user.IsBlocked {
		return nil, fmt.Errorf("USER_IS_BLOCKED")
	}

	userResponseDTO := userDTOs.UserResponseDTO{
		Id:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		Role:      user.Role.Name,
		IsBlocked: user.IsBlocked,
	}

	return &userResponseDTO, nil
}

func (s *authService) ChangePassword(
	userId uint,
	changePasswordDTO *auth.ChangePasswordDTO,
) error {
	user, err := s.userRepo.GetById(s.context, userId)
	if err != nil {
		return err
	}

	if !password_utils.ComparePasswords(user.Password, changePasswordDTO.OldPassword) {
		return fmt.Errorf("INCORRECT_USERNAME_OR_PASSWORD")
	}

	hashedPassword, err := password_utils.HashPassword(changePasswordDTO.NewPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	s.userRepo.Save(s.context, user)

	return nil
}

func (s *authService) RequestPasswordReset(
	requestPasswordResetDTO *auth.RequestPasswordResetDTO,
) error {
	return s.trManager.Do(s.context, func(ctx context.Context) error {
		user, err := s.userRepo.GetByEmail(ctx, requestPasswordResetDTO.Email)
		if err != nil {
			return err
		}

		// Get existing password recovery token
		passwordRecoveryToken, err := s.recoveryRepo.GetByUserId(ctx, user.ID)
		// If not found, create new one
		if err != nil {
			passwordRecoveryToken = &models.PasswordRecoveryToken{
				Token: password_utils.GenerateToken(),
				User:  *user,
			}

			err = s.recoveryRepo.Create(ctx, passwordRecoveryToken)
			if err != nil {
				return err
			}
		}

		err = s.emailService.SendPasswordResetEmail(
			user,
			passwordRecoveryToken.Token,
		)

		return err
	})
}

func (s *authService) ResetPassword(
	resetPasswordDTO *auth.ResetPasswordDTO,
) error {
	return s.trManager.Do(s.context, func(ctx context.Context) error {

		passwordRecoveryToken, err := s.recoveryRepo.GetByToken(ctx, resetPasswordDTO.Token)
		if err != nil {
			return err
		}

		user := passwordRecoveryToken.User
		hashedPassword, err := password_utils.HashPassword(resetPasswordDTO.Password)
		if err != nil {
			return err
		}

		user.Password = hashedPassword
		err = s.userRepo.Save(ctx, &user)
		if err != nil {
			return err
		}

		err = s.recoveryRepo.Delete(ctx, passwordRecoveryToken)
		if err != nil {
			return err
		}

		return nil
	})
}
