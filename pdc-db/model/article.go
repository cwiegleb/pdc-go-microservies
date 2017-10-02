package model

import (
	"github.com/jinzhu/gorm"
)

// Article Model
type Article struct {
	gorm.Model
	Text      string
	Size      string
	DealerID  uint
	Available bool
	Costs     float64
	Currency  string
}
