package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/cwiegleb/pdc-services/pdc-db/config"
	"github.com/cwiegleb/pdc-services/pdc-db/model"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type AccountingResult struct {
	Count_orders int
	Total_amount float32
}

func GetAccountingHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var config = config.LoadConfiguration("")
	db, err := gorm.Open(config.DbDriver, config.DbConnection)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print("failed to connect database", err)
		return
	}
	defer db.Close()

	var cashboxGet model.Cashbox
	if db.Where("ID = ?", vars["id"]).First(&cashboxGet).Error != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("Entry %s not found", vars["id"])
		return
	}

	var accountingResult AccountingResult
	db.Raw("SELECT count(orders.id) AS count_orders, sum(order_lines.price) AS total_amount FROM cashboxes, orders, order_lines WHERE cashboxes.id = orders.cashbox_id AND orders.id = order_lines.order_id AND orders.updated_at BETWEEN cashboxes.valid_from_date AND cashboxes.valid_to_date AND cashboxes.id = ?", vars["id"]).Scan(&accountingResult)

	b, err := json.Marshal(accountingResult)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(b)
}
