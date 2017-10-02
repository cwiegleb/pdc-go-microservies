package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/cwiegleb/pdc-services/pdc-db/config"
	"github.com/cwiegleb/pdc-services/pdc-db/model"
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
	db.AutoMigrate(&model.Article{})

	defer db.Close()
}
