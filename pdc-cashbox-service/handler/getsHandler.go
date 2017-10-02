package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/cwiegleb/pdc-services/pdc-db/config"
	"github.com/cwiegleb/pdc-services/pdc-db/model"
)

func GetsHandler(w http.ResponseWriter, r *http.Request) {

	var config = config.LoadConfiguration("")
	db, err := gorm.Open(config.DbDriver, config.DbConnection)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print("failed to connect database", err)
		return
	}
	defer db.Close()

	var cashboxesGet []model.Cashbox

	// @TODO check go date handling
	// var currentDate = Date.currentTime()
	//if db.Where("valid_from_date >= ? AND valid_to_date <= ?", currentDate, currentDate).Find(&cashboxesGet).Error != nil {
	if db.Find(&cashboxesGet).Error != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("No Entries found")
		return
	}

	b, err := json.Marshal(cashboxesGet)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(b)
}
