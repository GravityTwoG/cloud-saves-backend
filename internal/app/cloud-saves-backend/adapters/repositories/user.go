package repositories

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/adapters/models"
	"cloud-saves-backend/internal/app/cloud-saves-backend/domain/common"
	"cloud-saves-backend/internal/app/cloud-saves-backend/domain/user"
	gorm_utils "cloud-saves-backend/internal/app/cloud-saves-backend/utils/gorm"
	"context"

	trmgorm "github.com/avito-tech/go-transaction-manager/drivers/gorm/v2"
	"gorm.io/gorm"
)

// Implementing services.AuthRepository interface
type userRepo struct {
	db     *gorm.DB
	getter *trmgorm.CtxGetter
}

func NewUserRepository(db *gorm.DB, getter *trmgorm.CtxGetter) user.UserRepository {
	return &userRepo{
		db:     db,
		getter: getter,
	}
}

func (r *userRepo) Create(ctx context.Context, user *user.User) error {
	db := r.getter.DefaultTrOrDB(ctx, r.db)

	userModel := models.UserFromEntity(user)
	err := db.Preload("Role").Create(&userModel).Error
	if err != nil {
		return err
	}
	*user = *models.UserFromModel(userModel)
	return nil
}

func (r *userRepo) Save(ctx context.Context, user *user.User) error {
	db := r.getter.DefaultTrOrDB(ctx, r.db)

	userModel := models.UserFromEntity(user)
	err := db.Preload("Role").Save(userModel).Error
	if err != nil {
		return err
	}
	*user = *models.UserFromModel(userModel)
	return nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.db)

	userModel := models.User{}
	err := db.Where(&models.User{Email: email}).First(&userModel).Error
	if err != nil {
		return nil, err
	}

	return models.UserFromModel(&userModel), nil
}

func (r *userRepo) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.db)

	userModel := models.User{}
	err := db.Preload("Role").Where(
		&models.User{Username: username},
	).First(&userModel).Error
	if err != nil {
		return nil, err
	}

	return models.UserFromModel(&userModel), nil
}

func (r *userRepo) GetById(ctx context.Context, userId uint) (*user.User, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.db)

	userModel := models.User{}
	err := db.Preload("Role").First(&userModel, userId).Error
	if err != nil {
		return nil, err
	}

	return models.UserFromModel(&userModel), nil
}

func whereUsernameLike(searchQuery string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("username LIKE ?", "%"+searchQuery+"%")
	}
}

func (r *userRepo) GetUsersWithRole(
	ctx context.Context, 
	dto common.GetResourceDTO, 
	role *user.Role,
) (*common.ResourceDTO[user.User], error) {
	db := r.getter.DefaultTrOrDB(ctx, r.db)

	userModels := []models.User{}
	var totalCount int64

	err := db.
		Model(&models.User{}).
		Scopes(whereUsernameLike(dto.SearchQuery)).
		Where("role_id = ?", role.GetId()).
		Count(&totalCount).
		Order("username asc").
		Scopes(gorm_utils.Paginate(db, dto)).
		Preload("Role").
		Find(&userModels).Error
	if err != nil {
		return nil, err
	}

	users := make([]user.User, len(userModels))
	for i, userModel := range userModels {
		users[i] = *models.UserFromModel(&userModel)
	}

	return &common.ResourceDTO[user.User]{
		Items:      users,
		TotalCount: int(totalCount),
	}, nil
}
