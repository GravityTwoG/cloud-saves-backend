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
type roleRepo struct {
	db     *gorm.DB
	getter *trmgorm.CtxGetter
}

func NewRoleRepository(db *gorm.DB, getter *trmgorm.CtxGetter) auth.RoleRepository {
	return &roleRepo{
		db:     db,
		getter: getter,
	}
}

func (r *roleRepo) GetByName(ctx context.Context, name user.RoleName) (*user.Role, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.db)

	roleModel := models.Role{}
	err := db.Where(&models.Role{Name: name}).First(&roleModel).Error
	if err != nil {
		return nil, err
	}

	return models.RoleFromModel(&roleModel), nil
}
