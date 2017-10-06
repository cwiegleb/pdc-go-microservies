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
		var article model.Article
		if db.Where("id = ? AND available = 0", oldOrderLines[i].ArticleID).First(&article).Error != nil {
			tx.Rollback()
			w.WriteHeader(http.StatusBadRequest)
			log.Print("cannot find old article")
			return
		}
		article.Available = true
		if tx.Save(&article).Error != nil {
			tx.Rollback()
			w.WriteHeader(http.StatusBadRequest)
			log.Print("failed to update old article")
			return
		}

		if tx.Delete(&model.OrderLine{}, oldOrderLines[i].ID).Error != nil {
			tx.Rollback()
			w.WriteHeader(http.StatusBadRequest)
			log.Print("failed to delete old orderlines")
			return
		}
	}

	for i := 0; i < len(order.OrderLines); i++ {
		var article model.Article

		if order.OrderLines[i].ArticleID != 0 {
			if db.Where("id = ? AND available = 1", order.OrderLines[i].ArticleID).First(&article).Error != nil {
				tx.Rollback()
				w.WriteHeader(http.StatusBadRequest)
				log.Print("article already sold")
				return
			}

			var articleUpdate1 model.Article
			if tx.Model(&articleUpdate1).Where("id = ? AND available = 1", order.OrderLines[i].ArticleID).Update("available", 0).Error != nil {
				tx.Rollback()
				w.WriteHeader(http.StatusBadRequest)
				log.Print("failed to update article")
				return
			}
		}

		tx.NewRecord(order.OrderLines[i])
	}

	tx.Commit()
	defer tx.Close()
}
