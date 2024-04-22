package auth

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/adapters/services"
	"cloud-saves-backend/internal/app/cloud-saves-backend/domain/user"
	"cloud-saves-backend/internal/app/cloud-saves-backend/ports/dto/auth"
	"context"
	"fmt"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
)

type AuthService interface {
	Register(ctx context.Context, dto *auth.RegisterDTO) (*user.User, error)

	Login(ctx context.Context, dto *auth.LoginDTO) (*user.User, error)

	ChangePassword(ctx context.Context, userId uint, dto *auth.ChangePasswordDTO) error

	RequestPasswordReset(ctx context.Context, dto *auth.RequestPasswordResetDTO) error

	ResetPassword(ctx context.Context, dto *auth.ResetPasswordDTO) error

	BlockUser(ctx context.Context, userId uint) error

	UnblockUser(ctx context.Context, userId uint) error
}

type authService struct {
	trManager trm.Manager

	roleRepo     user.RoleRepository
	userRepo     user.UserRepository
	recoveryRepo PasswordRecoveryTokenRepository

	emailService services.EmailService
}

func NewAuth(
	trManager trm.Manager,
	roleRepo user.RoleRepository,
	userRepo user.UserRepository,
	recoveryRepo PasswordRecoveryTokenRepository,
	emailService services.EmailService,
) AuthService {
	return &authService{
		trManager:    trManager,
		roleRepo:     roleRepo,
		userRepo:     userRepo,
		recoveryRepo: recoveryRepo,
		emailService: emailService,
	}
}

func (s *authService) Register(ctx context.Context, registerDTO *auth.RegisterDTO) (*user.User, error) {
	roleUser, err := s.roleRepo.GetByName(ctx, user.RoleUser)
	if err != nil {
		return nil, err
	}

	user, err := user.NewUser(
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

	return user, nil
}

func (s *authService) Login(
	ctx context.Context,
	loginDTO *auth.LoginDTO,
) (*user.User, error) {
	user, err := s.userRepo.GetByUsername(ctx, loginDTO.Username)
	if err != nil {
		return nil, err
	}

	if !user.ComparePassword(loginDTO.Password) {
		return nil, fmt.Errorf("INCORRECT_USERNAME_OR_PASSWORD")
	}
	if user.IsBlocked() {
		return nil, fmt.Errorf("USER_IS_BLOCKED")
	}

	return user, nil
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

	err = user.ChangePassword(changePasswordDTO.NewPassword)
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
		recoveryToken, err := s.recoveryRepo.GetByUserId(ctx, user.GetId())
		// If not found, create new one
		if err != nil {
			recoveryToken = NewPasswordRecoveryToken(user)
			err = s.recoveryRepo.Create(ctx, recoveryToken)
			if err != nil {
				return err
			}
		}

		return s.emailService.SendPasswordResetEmail(
			ctx,
			user,
			recoveryToken.GetToken(),
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

		user := passwordRecoveryToken.GetUser()
		err = user.ChangePassword(resetPasswordDTO.Password)
		if err != nil {
			return err
		}

		err = s.userRepo.Save(ctx, user)
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

	user.Block()

	return s.userRepo.Save(ctx, user)
}

func (s *authService) UnblockUser(ctx context.Context, userId uint) error {
	user, err := s.userRepo.GetById(ctx, userId)
	if err != nil {
		return err
	}

	user.Unblock()

	return s.userRepo.Save(ctx, user)
}
