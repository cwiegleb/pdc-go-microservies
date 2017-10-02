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

	var articlesGet []model.Article
	if db.Where("dealer_id = ? AND available = 1", vars["id"]).Find(&articlesGet).Error != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("Entries %s not found", vars["id"])
		return
	}
	b, err := json.Marshal(articlesGet)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(b)
}
