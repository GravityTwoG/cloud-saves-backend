package repositories

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/models"
	"cloud-saves-backend/internal/app/cloud-saves-backend/services"
	"context"

	trmgorm "github.com/avito-tech/go-transaction-manager/drivers/gorm/v2"
	"gorm.io/gorm"
)

// Implementing services.RoleRepository interface
type roleRepo struct {
	db     *gorm.DB
	getter *trmgorm.CtxGetter
}

func NewRoleRepository(db *gorm.DB, getter *trmgorm.CtxGetter) services.RoleRepository {
	return &roleRepo{
		db:     db,
		getter: getter,
	}
}

func (r *roleRepo) GetByName(ctx context.Context, name models.RoleName) (*models.Role, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.db)

	role := models.Role{}
	err := db.Where(&models.Role{Name: name}).First(&role).Error
	if err != nil {
		return nil, err
	}

	return &role, nil
}
