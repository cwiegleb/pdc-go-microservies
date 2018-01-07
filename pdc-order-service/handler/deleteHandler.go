package handler

import (
	"log"
	"net/http"

	"github.com/cwiegleb/pdc-services/pdc-db/config"
	"github.com/cwiegleb/pdc-services/pdc-db/model"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var config = config.LoadConfiguration("")
	db, err := gorm.Open(config.DbDriver, config.DbConnection)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print("failed to connect database", err)
		return
	}
	defer db.Close()

	if db.Delete(&model.Order{}, vars["order-id"]).Error != nil {
		log.Printf("Entry %s not found", vars["order-id"])
		w.WriteHeader(http.StatusNotFound)
		return
	}

	tx := db.Begin()

	if tx.Delete(&model.Order{}, vars["order-id"]).Error != nil {
		tx.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		log.Print("failed to delete order")
		return
	}

	var oldOrderLines []model.OrderLine
	if db.Where("order_id = ?", vars["order-id"]).Find(&oldOrderLines).Error != nil {
		tx.Rollback()
		w.WriteHeader(http.StatusBadRequest)
		log.Print("cannot read old orderlines")
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
	tx.Commit()
	defer tx.Close()
}
