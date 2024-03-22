package gorm_utils

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/dto/common"

	"gorm.io/gorm"
)

func Paginate(db *gorm.DB, dto common.GetResourceDTO) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset((dto.PageNumber - 1) * dto.PageSize).Limit(dto.PageSize)
	}
}
