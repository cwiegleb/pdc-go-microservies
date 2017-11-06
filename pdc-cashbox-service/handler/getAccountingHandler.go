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

	var accountingResult []model.CashboxAccounting
	db.Raw("select to_char(orders.updated_at at time zone 'CET', 'DD.MM.YYYY HH24 Uhr') as order_date, count(DISTINCT orders.id) as count_orders, sum(order_lines.price) as total_amount from cashboxes, orders, order_lines where cashboxes.id = orders.cashbox_id and orders.id = order_lines.order_id and orders.updated_at between cashboxes.valid_from_date and cashboxes.valid_to_date and cashboxes.id = ? GROUP BY order_date, trunc(EXTRACT(hour from orders.updated_at) / 24)", vars["id"]).Scan(&accountingResult)

	b, err := json.Marshal(accountingResult)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(b)
}
