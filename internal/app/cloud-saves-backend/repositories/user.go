package repositories

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/dto/common"
	"cloud-saves-backend/internal/app/cloud-saves-backend/models"
	"cloud-saves-backend/internal/app/cloud-saves-backend/services"
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

func NewUserRepository(db *gorm.DB, getter *trmgorm.CtxGetter) services.UserRepository {
	return &userRepo{
		db:     db,
		getter: getter,
	}
}

func (r *userRepo) Create(ctx context.Context, user *models.User) error {
	db := r.getter.DefaultTrOrDB(ctx, r.db)

	return db.Preload("Role").Create(&user).Error
}

func (r *userRepo) Save(ctx context.Context, user *models.User) error {
	db := r.getter.DefaultTrOrDB(ctx, r.db)

	return db.Preload("Role").Save(user).Error
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.db)

	user := models.User{}
	err := db.Where(&models.User{Email: email}).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.db)

	user := models.User{}
	err := db.Preload("Role").Where(
		&models.User{Username: username},
	).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) GetById(ctx context.Context, userId uint) (*models.User, error) {
	db := r.getter.DefaultTrOrDB(ctx, r.db)

	user := models.User{}
	err := db.Preload("Role").First(&user, userId).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func whereUsernameLike(searchQuery string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("username LIKE ?", "%"+searchQuery+"%")
	}
}

func (r *userRepo) GetAll(ctx context.Context, dto common.GetResourceDTO) (*common.ResourceDTO[models.User], error) {
	db := r.getter.DefaultTrOrDB(ctx, r.db)

	users := []models.User{}
	var totalCount int64

	err := db.
		Model(&models.User{}).
		Scopes(whereUsernameLike(dto.SearchQuery)).
		Count(&totalCount).
		Order("username asc").
		Scopes(gorm_utils.Paginate(db, dto)).
		Preload("Role").
		Find(&users).Error
	if err != nil {
		return nil, err
	}

	return &common.ResourceDTO[models.User]{
		Items:      users,
		TotalCount: int(totalCount),
	}, nil
}
