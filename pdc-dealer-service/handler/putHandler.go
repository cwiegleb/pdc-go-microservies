package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/cwiegleb/pdc-services/pdc-db/config"
	"github.com/cwiegleb/pdc-services/pdc-db/model"
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

	dealer := &model.Dealer{}
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if json.Unmarshal(bodyBytes, dealer) != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		log.Print("Cannot unmarshal dealer ")
		return
	}

	var dealerGet model.Dealer
	if db.Where(vars["id"]).First(&dealerGet).Error != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("Entry %s not found", vars["id"])
		return
	}

	if strconv.Itoa(int(dealer.ID)) != vars["id"] {
		w.WriteHeader(http.StatusBadRequest)
		log.Print("IDs mismatched")
		return
	}

	if err := db.Save(dealer).Error; err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		log.Print("Cannot save dealer ", err)
	}
}
