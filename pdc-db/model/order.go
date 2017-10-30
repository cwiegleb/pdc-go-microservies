package model

import (
	"github.com/jinzhu/gorm"
)

// Order Model
type Order struct {
	gorm.Model
	CashboxID   uint
	OrderStatus uint
	OrderLines  []OrderLine
}
