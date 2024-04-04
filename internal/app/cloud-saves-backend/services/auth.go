package services

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/dto/auth"
	userDTOs "cloud-saves-backend/internal/app/cloud-saves-backend/dto/user"
	"cloud-saves-backend/internal/app/cloud-saves-backend/models"
	"context"
	"fmt"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
)

type RoleRepository interface {
	GetByName(ctx context.Context, name models.RoleName) (*models.Role, error)
}

type PasswordRecoveryTokenRepository interface {
	Create(ctx context.Context, token *models.PasswordRecoveryToken) error

	Save(ctx context.Context, token *models.PasswordRecoveryToken) error

	GetByToken(ctx context.Context, token string) (*models.PasswordRecoveryToken, error)

	GetByUserId(ctx context.Context, userId uint) (*models.PasswordRecoveryToken, error)

	Delete(ctx context.Context, token *models.PasswordRecoveryToken) error
}

type AuthService interface {
	Register(ctx context.Context, dto *auth.RegisterDTO) (*userDTOs.UserResponseDTO, error)

	Login(ctx context.Context, dto *auth.LoginDTO) (*userDTOs.UserResponseDTO, error)

	ChangePassword(ctx context.Context, userId uint, dto *auth.ChangePasswordDTO) error

	RequestPasswordReset(ctx context.Context, dto *auth.RequestPasswordResetDTO) error

	ResetPassword(ctx context.Context, dto *auth.ResetPasswordDTO) error

	BlockUser(ctx context.Context, userId uint) error

	UnblockUser(ctx context.Context, userId uint) error
}

type authService struct {
	trManager trm.Manager

	roleRepo     RoleRepository
	userRepo     UserRepository
	recoveryRepo PasswordRecoveryTokenRepository

	emailService EmailService
}

func NewAuth(
	trManager trm.Manager,
	roleRepo RoleRepository,
	userRepo UserRepository,
	recoveryRepo PasswordRecoveryTokenRepository,
	emailService EmailService,
) AuthService {
	return &authService{
		trManager:    trManager,
		roleRepo:     roleRepo,
		userRepo:     userRepo,
		recoveryRepo: recoveryRepo,
		emailService: emailService,
	}
}

func (s *authService) Register(ctx context.Context, registerDTO *auth.RegisterDTO) (*userDTOs.UserResponseDTO, error) {
	roleUser, err := s.roleRepo.GetByName(ctx, models.RoleUser)
	if err != nil {
		return nil, err
	}

	user, err := models.NewUser(
		registerDTO.Username,
		registerDTO.Email,
		registerDTO.Password,
		roleUser,
	)
	if err != nil {
		return nil, err
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return userDTOs.FromUser(user), nil
}

func (s *authService) Login(
	ctx context.Context,
	loginDTO *auth.LoginDTO,
) (*userDTOs.UserResponseDTO, error) {
	user, err := s.userRepo.GetByUsername(ctx, loginDTO.Username)
	if err != nil {
		return nil, err
	}

	if !user.ComparePassword(loginDTO.Password) {
		return nil, fmt.Errorf("INCORRECT_USERNAME_OR_PASSWORD")
	}
	if user.IsBlocked {
		return nil, fmt.Errorf("USER_IS_BLOCKED")
	}

	return userDTOs.FromUser(user), nil
}

func (s *authService) ChangePassword(
	ctx context.Context,
	userId uint,
	changePasswordDTO *auth.ChangePasswordDTO,
) error {
	user, err := s.userRepo.GetById(ctx, userId)
	if err != nil {
		return err
	}

	if !user.ComparePassword(changePasswordDTO.OldPassword) {
		return fmt.Errorf("INCORRECT_USERNAME_OR_PASSWORD")
	}

	err = user.SetPassword(changePasswordDTO.NewPassword)
	if err != nil {
		return err
	}

	return s.userRepo.Save(ctx, user)
}

func (s *authService) RequestPasswordReset(
	ctx context.Context,
	requestPasswordResetDTO *auth.RequestPasswordResetDTO,
) error {
	return s.trManager.Do(ctx, func(ctx context.Context) error {
		user, err := s.userRepo.GetByEmail(ctx, requestPasswordResetDTO.Email)
		if err != nil {
			return err
		}

		// Get existing password recovery token
		recoveryToken, err := s.recoveryRepo.GetByUserId(ctx, user.ID)
		// If not found, create new one
		if err != nil {
			recoveryToken = models.NewPasswordRecoveryToken(user)
			err = s.recoveryRepo.Create(ctx, recoveryToken)
			if err != nil {
				return err
			}
		}

		return s.emailService.SendPasswordResetEmail(
			ctx,
			user,
			recoveryToken.Token,
		)
	})
}

func (s *authService) ResetPassword(
	ctx context.Context,
	resetPasswordDTO *auth.ResetPasswordDTO,
) error {
	return s.trManager.Do(ctx, func(ctx context.Context) error {
		passwordRecoveryToken, err := s.recoveryRepo.GetByToken(ctx, resetPasswordDTO.Token)
		if err != nil {
			return err
		}

		user := passwordRecoveryToken.User
		err = user.SetPassword(resetPasswordDTO.Password)
		if err != nil {
			return err
		}

		err = s.userRepo.Save(ctx, &user)
		if err != nil {
			return err
		}

		return s.recoveryRepo.Delete(ctx, passwordRecoveryToken)
	})
}

func (s *authService) BlockUser(ctx context.Context, userId uint) error {
	user, err := s.userRepo.GetById(ctx, userId)
	if err != nil {
		return err
	}

	user.IsBlocked = true

	return s.userRepo.Save(ctx, user)
}

func (s *authService) UnblockUser(ctx context.Context, userId uint) error {
	user, err := s.userRepo.GetById(ctx, userId)
	if err != nil {
		return err
	}

	user.IsBlocked = false

	return s.userRepo.Save(ctx, user)
}
