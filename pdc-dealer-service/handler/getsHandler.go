package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/cwiegleb/pdc-services/pdc-db/config"
	"github.com/cwiegleb/pdc-services/pdc-db/model"
	"github.com/jinzhu/gorm"
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

	var dealersGet []model.Dealer
	if db.Find(&dealersGet).Error != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("No Entries found")
		return
	}
	b, err := json.Marshal(dealersGet)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}
