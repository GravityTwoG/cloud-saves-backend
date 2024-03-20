package services

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/dto/auth"
	userDTOs "cloud-saves-backend/internal/app/cloud-saves-backend/dto/user"
	"cloud-saves-backend/internal/app/cloud-saves-backend/models"
	password_utils "cloud-saves-backend/internal/app/cloud-saves-backend/utils/password-utils"
	"fmt"

	"gorm.io/gorm"
)

type AuthService interface {
	Register(dto *auth.RegisterDTO) (*userDTOs.UserResponseDTO, error)
	
	Login(dto  *auth.LoginDTO)  (*userDTOs.UserResponseDTO, error)
	
	ChangePassword(userId uint, dto *auth.ChangePasswordDTO) error
	
	RequestPasswordReset(dto *auth.RequestPasswordResetDTO) error
	
	ResetPassword(dto *auth.ResetPasswordDTO) error
}

type authService struct {
	db *gorm.DB
	emailService EmailService
}

func NewAuth(db *gorm.DB, emailService EmailService) AuthService {
	return &authService{
		db: db,
		emailService: emailService,
	}
}

func (s *authService) Register(registerDTO *auth.RegisterDTO) (*userDTOs.UserResponseDTO, error) {
	roleUser := models.Role{} 
	err := s.db.Where(&models.Role{Name: "ROLE_USER"}).First(&roleUser).Error
	if err != nil {
		return nil, err
	}
	
	hashedPassword, err := password_utils.HashPassword(registerDTO.Password);
	if err != nil {
		return nil, err
	}
	
	user := models.User{
		Username: registerDTO.Username,
		Email: registerDTO.Email,
		Password: hashedPassword,
		IsBlocked: false,
		Role: roleUser,
	}
	
	err = s.db.Preload("Role").Create(&user).Error
	if err != nil {
		return nil, err
	}
	
	userResponseDTO := userDTOs.UserResponseDTO{
		Id: user.ID,
		Email: user.Email,
		Username: user.Username,
		Role: user.Role.Name,
		IsBlocked: user.IsBlocked,
	}
	
	return &userResponseDTO, nil
}

func (s *authService) Login(
	loginDTO *auth.LoginDTO,
) (*userDTOs.UserResponseDTO, error) {
	user := models.User{};

	err := s.db.Preload("Role").Where(
		&models.User{Username: loginDTO.Username},
	).First(&user).Error
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
		Id: user.ID,
		Email: user.Email,
		Username: user.Username,
		Role: user.Role.Name,
		IsBlocked: user.IsBlocked,
	}

	return &userResponseDTO, nil
}

func (s *authService) ChangePassword(
	userId uint,
	changePasswordDTO *auth.ChangePasswordDTO,
) error {
	user := models.User{};

	err := s.db.First(&user, userId).Error
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
	s.db.Save(&user)

	return nil
}

func (s *authService) RequestPasswordReset(
	requestPasswordResetDTO *auth.RequestPasswordResetDTO,
) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		filter := models.User{Email: requestPasswordResetDTO.Email};
		user := models.User{};
		err := tx.Where(&filter).First(&user).Error
		if err != nil {
			return err
		}

		passwordRecoveryToken := models.PasswordRecoveryToken{
			Token: password_utils.GenerateToken(),
			User: user,
		}
	
		err = tx.Create(&passwordRecoveryToken).Error
		if err != nil {
			return err
		}
	
		err = s.emailService.SendPasswordResetEmail(
			&user, 
			passwordRecoveryToken.Token,
		)
	
		return err
	})
}

func (s *authService) ResetPassword(
	resetPasswordDTO *auth.ResetPasswordDTO,
) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		filter := models.PasswordRecoveryToken{Token: resetPasswordDTO.Token};
		passwordRecoveryToken := models.PasswordRecoveryToken{};
		err := tx.Where(&filter).Preload("User").First(&passwordRecoveryToken).Error
		if err != nil {
			return err
		}
		
		user := passwordRecoveryToken.User
		hashedPassword, err := password_utils.HashPassword(resetPasswordDTO.Password)
		if err != nil {
			return err
		}
		
		user.Password = hashedPassword
		err = tx.Save(&user).Error
		if err != nil {
			return err
		}
		
		err = tx.Delete(&passwordRecoveryToken).Error
		if err != nil {
			return err	
		}
		
		return nil
	})
}

