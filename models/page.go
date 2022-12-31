package models

import (
	"gorm.io/gorm"
)

type Page struct {
	PageNumber    int   //現在のページ
	PageSize      int   //１ページあたりの要素数
	TotalElements int64 //合計の要素数
	TotalPages    int   //合計ページ数
	PaginateInfos []PaginateInfo
}

type PaginateInfo struct {
	PageNumber int    //ページ番号
	PathParam  string //ページのクエリ
	Info       string //追加情報
}

func Paginate(page Page) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page.PageNumber - 1) * page.PageSize
		return db.Offset(offset).Limit(page.PageSize)
	}
}
