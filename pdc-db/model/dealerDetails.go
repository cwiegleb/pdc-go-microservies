package model

import (
	"github.com/jinzhu/gorm"
)

// DealerDetails Model
type DealerDetails struct {
	gorm.Model
	DealerID   uint
 	Type       string
	Name       string
	Street     string
	City       string
	PostalCode string
	Telephone  string
	Email      string
	Iban       string
	Bic        string
	BankName   string
	Fee        float32
	Commission float32
	Currency   string
}
