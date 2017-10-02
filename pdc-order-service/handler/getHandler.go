package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/cwiegleb/pdc-services/pdc-db/config"
	"github.com/cwiegleb/pdc-services/pdc-db/model"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var config = config.LoadConfiguration("")
	db, err := gorm.Open(config.DbDriver, config.DbConnection)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print("failed to connect database", err)
		return
	}
	defer db.Close()

	var orderGet model.Order
	if db.Where("ID = ? AND cashbox_id = ?", vars["order-id"], vars["id"]).First(&orderGet).Error != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("Entry %s not found", vars["id"])
		return
	}

	var orderLineGet []model.OrderLine
	if db.Where("order_id = ?", orderGet.ID).First(&orderLineGet).Error != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("No orderlines found for orderId %s", vars["order-id"])
		return
	}

	orderGet.OrderLines = orderLineGet

	b, err := json.Marshal(orderGet)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(b)
}
