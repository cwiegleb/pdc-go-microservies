package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Cashbox Model
type Cashbox struct {
	gorm.Model
	Name          string
	ValidFromDate time.Time
	ValidToDate   time.Time
	Orders        []Order
}
