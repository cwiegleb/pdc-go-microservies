package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/cwiegleb/pdc-services/pdc-db/config"
	"github.com/cwiegleb/pdc-services/pdc-db/model"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func PutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var config = config.LoadConfiguration("")
	db, err := gorm.Open(config.DbDriver, config.DbConnection)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print("failed to connect database", err)
		return
	}
	defer db.Close()
	order := &model.Order{}
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if json.Unmarshal(bodyBytes, order) != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		log.Print("Cannot unmarshal order ")
		return
	}

	// Check CashboxId
	var cashboxGet model.Cashbox
	if db.Where(vars["id"]).First(&cashboxGet).Error != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("Dealer Entry %s not found", vars["id"])
		return
	}

	if strconv.Itoa(int(cashboxGet.ID)) != vars["id"] {
		w.WriteHeader(http.StatusBadRequest)
		log.Print("DealerIDs mismatched")
		return
	}

	// Check OrderId
	if strconv.Itoa(int(order.ID)) != vars["order-id"] {
		w.WriteHeader(http.StatusBadRequest)
		log.Print("OrderIDs mismatched")
		return
	}

	tx := db.Begin()

	if tx.Save(order).Error != nil {
		tx.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		log.Print("failed to save order")
		return
	}

	var oldOrderLines []model.OrderLine

	if db.Where("order_id = ?", order.ID).Find(&oldOrderLines).Error != nil {
		tx.Rollback()
		w.WriteHeader(http.StatusBadRequest)
		log.Print("cannot read old orders")
		return
	}

	for i := 0; i < len(oldOrderLines); i++ {

		if tx.Delete(&model.OrderLine{}, oldOrderLines[i].ID).Error != nil {
			tx.Rollback()
			w.WriteHeader(http.StatusBadRequest)
			log.Print("failed to delete old orderlines")
			return
		}
	}

	for i := 0; i < len(order.OrderLines); i++ {
		tx.NewRecord(order.OrderLines[i])
	}

	tx.Commit()
	defer tx.Close()
}
