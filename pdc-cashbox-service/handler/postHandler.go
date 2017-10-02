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

	cashbox := &model.Cashbox{}
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if json.Unmarshal(bodyBytes, cashbox) != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		log.Print("Cannot unmarshal cashbox ")
		return
	}

	if db.Create(cashbox).Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print("failed to create cashbox", err)
		return
	}

	location := []string{r.Host, "cashboxes", strconv.Itoa(int(cashbox.ID))}
	w.Header().Set("Location", strings.Join(location, "/"))
}
