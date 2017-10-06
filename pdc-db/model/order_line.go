package model

import (
	"github.com/jinzhu/gorm"
)

// OrderLine Model
type OrderLine struct {
	gorm.Model
	OrderID       uint
	ArticleID     uint
	UnkownArticle string
	Price         float64
	Currency      string
}
