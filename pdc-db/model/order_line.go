package model

import (
	"github.com/jinzhu/gorm"
)

// OrderLine Model
type OrderLine struct {
	gorm.Model
	OrderID    uint
	ArticleID  uint
	DealerID   uint
	DealerText string
	Price      float64
	Currency   string
}
