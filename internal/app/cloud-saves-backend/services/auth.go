package services

import (
	authDTOs "cloud-saves-backend/internal/app/cloud-saves-backend/dto/auth"
	userDTOs "cloud-saves-backend/internal/app/cloud-saves-backend/dto/user"
	"cloud-saves-backend/internal/app/cloud-saves-backend/models"
	password_utils "cloud-saves-backend/internal/app/cloud-saves-backend/utils/password-utils"
	"fmt"

	"gorm.io/gorm"
)

type AuthService interface {
	Register(dto *authDTOs.RegisterDTO) (*userDTOs.UserResponseDTO, error)
	
	Login(dto  *authDTOs.LoginDTO)  (*userDTOs.UserResponseDTO, error)
	
	ChangePassword(userId uint, dto *authDTOs.ChangePasswordDTO) error
	
	RequestPasswordReset(userId uint, dto *authDTOs.RequestPasswordResetDTO) error
	
	ResetPassword(userId uint, dto *authDTOs.ResetPasswordDTO) error
}

type authService struct {
	db *gorm.DB
}

func NewAuth(db *gorm.DB) AuthService {
	return &authService{
		db: db,
	}
}

func (s *authService) Register(registerDTO *authDTOs.RegisterDTO) (*userDTOs.UserResponseDTO, error) {
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
	loginDTO *authDTOs.LoginDTO,
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
	changePasswordDTO *authDTOs.ChangePasswordDTO,
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
	userId uint,
	requestPasswordResetDTO *authDTOs.RequestPasswordResetDTO,
) error {
	panic("unimplemented")
}

func (s *authService) ResetPassword(
	userId uint,
	resetPasswordDTO *authDTOs.ResetPasswordDTO,
) error {
	panic("unimplemented")
}

