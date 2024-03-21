package repositories

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/models"
	"cloud-saves-backend/internal/app/cloud-saves-backend/services"
	"context"

	trmgorm "github.com/avito-tech/go-transaction-manager/drivers/gorm/v2"
	"gorm.io/gorm"
)

// Implementing services.RoleRepository interface
type recoveryRepo struct {
	db     *gorm.DB
	getter *trmgorm.CtxGetter
}

func NewPasswordRecoveryTokenRepository(db *gorm.DB, getter *trmgorm.CtxGetter) services.PasswordRecoveryTokenRepository {
	return &recoveryRepo{
		db:     db,
		getter: getter,
	}
}

func (r *recoveryRepo) Create(ctx context.Context, token *models.PasswordRecoveryToken) error {
	db := r.getter.DefaultTrOrDB(ctx, r.db)

	return db.Preload("User").Create(token).Error
}

func (r *recoveryRepo) Save(ctx context.Context, token *models.PasswordRecoveryToken) error {
	db := r.getter.DefaultTrOrDB(ctx, r.db)

	return db.Preload("User").Save(token).Error
}

func (r *recoveryRepo) GetByToken(ctx context.Context, token string) (*models.PasswordRecoveryToken, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.db)

	tokenModel := models.PasswordRecoveryToken{}
	err := db.Preload("User").Where(&models.PasswordRecoveryToken{Token: token}).First(&tokenModel).Error
	if err != nil {
		return nil, err
	}

	return &tokenModel, nil
}

func (r *recoveryRepo) GetByUserId(ctx context.Context, userId uint) (*models.PasswordRecoveryToken, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.db)

	token := models.PasswordRecoveryToken{}
	err := db.
		Where(&models.PasswordRecoveryToken{UserID: userId}).
		First(&token).Error
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *recoveryRepo) Delete(ctx context.Context, token *models.PasswordRecoveryToken) error {
	db := r.getter.DefaultTrOrDB(ctx, r.db)

	return db.Delete(token).Error
}
