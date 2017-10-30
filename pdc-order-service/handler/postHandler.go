package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/cwiegleb/pdc-services/pdc-db/config"
	"github.com/cwiegleb/pdc-services/pdc-db/model"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
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
		log.Print("Cannot unmarshal order ", json.Unmarshal(bodyBytes, order))
		return
	}

	// Check CashboxId
	var cashboxGet model.Cashbox
	if db.Where(vars["id"]).First(&cashboxGet).Error != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("Cashbox Entry %s not found", vars["id"])
		return
	}

	if strconv.Itoa(int(cashboxGet.ID)) != vars["id"] {
		w.WriteHeader(http.StatusBadRequest)
		log.Print("IDs mismatched")
		return
	}

	tx := db.Begin()

	if tx.Create(order).Error != nil {
		tx.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		log.Print("failed to create order")
		return
	}

	for i := 0; i < len(order.OrderLines); i++ {
		var article model.Article

		if order.OrderLines[i].ArticleID != 0 {

			if db.Where("id = ? AND available = true", order.OrderLines[i].ArticleID).First(&article).Error != nil {
				tx.Rollback()
				w.WriteHeader(http.StatusBadRequest)
				log.Print("article already sold")
				return
			}

			var articleUpdate model.Article
			if order.OrderLines[i].ArticleID != 9999 && tx.Model(&articleUpdate).Where("id = ? AND available = true", order.OrderLines[i].ArticleID).Update("available", false).Error != nil {
				tx.Rollback()
				w.WriteHeader(http.StatusBadRequest)
				log.Print("failed to update article")
				return
			}
		}

		order.OrderLines[i].OrderID = order.ID
		tx.NewRecord(order.OrderLines[i])
	}

	tx.Commit()
	defer tx.Close()

	location := []string{r.Host, "cashboxes", strconv.Itoa(int(cashboxGet.ID)), "orders", strconv.Itoa(int(order.ID))}
	w.Header().Set("Location", strings.Join(location, "/"))
}
