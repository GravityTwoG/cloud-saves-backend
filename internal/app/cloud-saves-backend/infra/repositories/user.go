package repositories

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/domain/user"
	"cloud-saves-backend/internal/app/cloud-saves-backend/dto/common"
	"cloud-saves-backend/internal/app/cloud-saves-backend/infra/models"
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
	return db.Preload("Role").Create(&userModel).Error
}

func (r *userRepo) Save(ctx context.Context, user *user.User) error {
	db := r.getter.DefaultTrOrDB(ctx, r.db)

	userModel := models.UserFromEntity(user)
	return db.Preload("Role").Save(userModel).Error
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.db)

	userModel := models.UserModel{}
	err := db.Where(&models.UserModel{Email: email}).First(&userModel).Error
	if err != nil {
		return nil, err
	}

	return models.UserFromModel(&userModel), nil
}

func (r *userRepo) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.db)

	userModel := models.UserModel{}
	err := db.Preload("Role").Where(
		&models.UserModel{Username: username},
	).First(&userModel).Error
	if err != nil {
		return nil, err
	}

	return models.UserFromModel(&userModel), nil
}

func (r *userRepo) GetById(ctx context.Context, userId uint) (*user.User, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.db)

	userModel := models.UserModel{}
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

func (r *userRepo) GetAll(ctx context.Context, dto common.GetResourceDTO) (*common.ResourceDTO[user.User], error) {
	db := r.getter.DefaultTrOrDB(ctx, r.db)

	userModels := []models.UserModel{}
	var totalCount int64

	err := db.
		Model(&models.UserModel{}).
		Scopes(whereUsernameLike(dto.SearchQuery)).
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
