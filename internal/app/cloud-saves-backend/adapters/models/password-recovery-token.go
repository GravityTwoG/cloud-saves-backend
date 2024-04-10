package models

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/domain/auth"
	"time"
)

type PasswordRecoveryToken struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Token     string `gorm:"unique;not null"`
	UserID    uint   `gorm:"unique;not null"`
	User      User   `gorm:"onDelete:CASCADE"`
}

func PasswordRecoveryTokenFromEntity(token *auth.PasswordRecoveryToken) *PasswordRecoveryToken {
	return &PasswordRecoveryToken{
		ID:        token.GetId(),
		CreatedAt: token.GetCreatedAt(),
		UpdatedAt: token.GetUpdatedAt(),
		Token:     token.GetToken(),
		UserID:    token.GetUser().GetId(),
		User:      *UserFromEntity(token.GetUser()),
	}
}

func PasswordRecoveryTokenFromModel(tokenModel *PasswordRecoveryToken) *auth.PasswordRecoveryToken {
	return auth.PasswordRecoveryTokenFromDB(
		tokenModel.ID,
		tokenModel.Token,
		UserFromModel(&tokenModel.User),
		tokenModel.CreatedAt,
		tokenModel.UpdatedAt,
	)
}
