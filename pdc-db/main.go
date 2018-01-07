package main

import (
	"github.com/cwiegleb/pdc-services/pdc-db/config"
	"github.com/cwiegleb/pdc-services/pdc-db/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {

	var config config.Config = config.LoadConfiguration("")
	db, err := gorm.Open(config.DbDriver, config.DbConnection)
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&model.Cashbox{})
	db.AutoMigrate(&model.Order{})
	db.AutoMigrate(&model.OrderLine{})
	db.AutoMigrate(&model.Dealer{})
	db.AutoMigrate(&model.DealerDetails{})

	var unkwonDealer = &model.Dealer{
		Text:       "Unbekannter Anbieter",
		ExternalId: "Unbekannter Anbieter",
	}
	unkwonDealer.ID = 9999

	if db.Create(unkwonDealer).Error != nil {
		println(db.Create(unkwonDealer).Error)
	}

	defer db.Close()
}
