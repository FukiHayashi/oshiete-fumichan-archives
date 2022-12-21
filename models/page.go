package models

import (
	"gorm.io/gorm"
)

type Page struct {
	PageNumber    int
	PageSize      int
	TotalElements int64
	TotalPages    int
	NextPage      int
	PrevPage      int
	PaginateNum   []int
}

func Paginate(page Page) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page.PageNumber - 1) * page.PageSize
		return db.Offset(offset).Limit(page.PageSize)
	}
}
