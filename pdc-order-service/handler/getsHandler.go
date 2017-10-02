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

func GetsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var config = config.LoadConfiguration("")
	db, err := gorm.Open(config.DbDriver, config.DbConnection)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print("failed to connect database", err)
		return
	}
	defer db.Close()

	var ordersGet []model.Order
	if db.Where("cashbox_id = ?", vars["id"]).Find(&ordersGet).Error != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("Entries %s not found", vars["id"])
		return
	}

	for i := 0; i < len(ordersGet); i++ {
		var orderLineGet []model.OrderLine
		if db.Where("order_id = ?", ordersGet[i].ID).First(&orderLineGet).Error != nil {
			w.WriteHeader(http.StatusNotFound)
			log.Printf("No orderlines found for orderId %s", vars["order-id"])
			return
		}
		ordersGet[i].OrderLines = orderLineGet
	}

	b, err := json.Marshal(ordersGet)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(b)
}
