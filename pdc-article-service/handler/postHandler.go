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
		log.Print("IDs mismatched")
		return
	}

	if db.Create(article).Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print("failed to create article", err)
		return
	}

	location := []string{r.Host, "articles", strconv.Itoa(int(article.ID)), "dealers", strconv.Itoa(int(dealerGet.ID))}
	w.Header().Set("Location", strings.Join(location, "/"))
}
