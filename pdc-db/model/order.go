package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Order Model
type Order struct {
	gorm.Model
	CashboxID    uint
	OrderStatus  string
	CreationDate time.Time
	OrderLines   []OrderLine
}
