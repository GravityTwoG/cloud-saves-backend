package common

type GetResourceDTO struct {
	PageNumber  int    `form:"pageNumber" binding:"required,gt=0"`
	PageSize    int    `form:"pageSize" binding:"required,gt=0"`
	SearchQuery string `form:"searchQuery" binding:"min=0"`
}
