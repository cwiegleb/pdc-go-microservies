package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/cwiegleb/pdc-services/pdc-db/config"
	"github.com/cwiegleb/pdc-services/pdc-db/model"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
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
	if db.Create(dealer).Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print("failed to create dealer", err)
		return
	}

	location := []string{r.Host, "dealers", strconv.Itoa(int(dealer.ID))}
	w.Header().Set("Location", strings.Join(location, "/"))
}
