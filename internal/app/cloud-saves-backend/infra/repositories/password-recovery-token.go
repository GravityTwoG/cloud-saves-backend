package repositories

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/domain/auth"
	"cloud-saves-backend/internal/app/cloud-saves-backend/domain/user"
	"cloud-saves-backend/internal/app/cloud-saves-backend/infra/models"
	"context"

	trmgorm "github.com/avito-tech/go-transaction-manager/drivers/gorm/v2"
	"gorm.io/gorm"
)

// Implementing services.RoleRepository interface
type recoveryRepo struct {
	db     *gorm.DB
	getter *trmgorm.CtxGetter
}

func NewPasswordRecoveryTokenRepository(db *gorm.DB, getter *trmgorm.CtxGetter) auth.PasswordRecoveryTokenRepository {
	return &recoveryRepo{
		db:     db,
		getter: getter,
	}
}

func (r *recoveryRepo) Create(ctx context.Context, token *auth.PasswordRecoveryToken) error {
	db := r.getter.DefaultTrOrDB(ctx, r.db)

	tokenModel := models.PasswordRecoveryTokenFromEntity(token)
	return db.Preload("User").Create(tokenModel).Error
}

func (r *recoveryRepo) Save(ctx context.Context, token *auth.PasswordRecoveryToken) error {
	db := r.getter.DefaultTrOrDB(ctx, r.db)

	tokenModel := models.PasswordRecoveryTokenFromEntity(token)
	return db.Preload("User").Save(tokenModel).Error
}

func (r *recoveryRepo) GetByToken(ctx context.Context, token string) (*auth.PasswordRecoveryToken, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.db)

	tokenModel := models.PasswordRecoveryToken{}
	err := db.Preload("User").Where(&models.PasswordRecoveryToken{Token: token}).First(&tokenModel).Error
	if err != nil {
		return nil, err
	}

	return auth.PasswordRecoveryTokenFromDB(
		tokenModel.ID,
		tokenModel.Token,
		user.UserFromDB(
			tokenModel.UserID,
			tokenModel.User.Username,
			tokenModel.User.Email,
			tokenModel.User.Password,
			user.RoleFromDB(tokenModel.User.Role.ID, tokenModel.User.Role.Name),
			tokenModel.User.IsBlocked,
		),
		tokenModel.CreatedAt,
		tokenModel.UpdatedAt,
	), nil
}

func (r *recoveryRepo) GetByUserId(ctx context.Context, userId uint) (*auth.PasswordRecoveryToken, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.db)

	tokenModel := models.PasswordRecoveryToken{}
	err := db.
		Where(&models.PasswordRecoveryToken{UserID: userId}).
		First(&tokenModel).Error
	if err != nil {
		return nil, err
	}

	return models.PasswordRecoveryTokenFromModel(&tokenModel), nil
}

func (r *recoveryRepo) Delete(ctx context.Context, token *auth.PasswordRecoveryToken) error {
	db := r.getter.DefaultTrOrDB(ctx, r.db)

	tokenModel := models.PasswordRecoveryTokenFromEntity(token)
	return db.Delete(tokenModel).Error
}
