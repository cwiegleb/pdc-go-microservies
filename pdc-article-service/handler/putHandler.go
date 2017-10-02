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
	article := &model.Article{}
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if json.Unmarshal(bodyBytes, article) != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		log.Print("Cannot unmarshal article ")
		return
	}

	// Check DealerId
	var dealerGet model.Dealer
	if db.Where(vars["id"]).First(&dealerGet).Error != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("Dealer Entry %s not found", vars["id"])
		return
	}

	if strconv.Itoa(int(dealerGet.ID)) != vars["id"] {
		w.WriteHeader(http.StatusBadRequest)
		log.Print("DealerIDs mismatched")
		return
	}

	// Check ArticleId
	if strconv.Itoa(int(article.ID)) != vars["article-id"] {
		w.WriteHeader(http.StatusBadRequest)
		log.Print("ArticleIDs mismatched")
		return
	}

	if db.Save(article).Error != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		log.Print("Cannot save customer ", err)
		return
	}
}
