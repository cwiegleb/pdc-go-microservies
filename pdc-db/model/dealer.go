package model

import (
	"github.com/jinzhu/gorm"
)

// Dealer Model
type Dealer struct {
	gorm.Model
	Text     string
	Articles []Article
}
