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
	cashbox := &model.Cashbox{}
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if json.Unmarshal(bodyBytes, cashbox) != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		log.Print("Cannot unmarshal cashbox ")
		return
	}

	// Check DealerId
	if strconv.Itoa(int(cashbox.ID)) != vars["id"] {
		w.WriteHeader(http.StatusBadRequest)
		log.Print("CashboxIDs mismatched")
		return
	}

	if db.Save(cashbox).Error != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		log.Print("Cannot save cashbox ", err)
		return
	}
}
