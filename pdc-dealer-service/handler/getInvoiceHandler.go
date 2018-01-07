package handler

import (
	"log"
	"net/http"

	"github.com/cwiegleb/pdc-services/pdc-db/config"
	"github.com/cwiegleb/pdc-services/pdc-db/model"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

/*
GetInvoiceHandler
*/
func GetInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var config = config.LoadConfiguration("")
	db, err := gorm.Open(config.DbDriver, config.DbConnection)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print("failed to connect database", err)
		return
	}
	defer db.Close()

	var dealerGet model.Dealer
	if db.Where("ID = ?", vars["id"]).First(&dealerGet).Error != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("Entry %s not found", vars["id"])
		return
	}

	var dealerDetails model.DealerDetails
	if db.Where("dealer_id = ?", vars["id"]).First(&dealerDetails).Error != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("Details entry %s not found", vars["id"])
		return
	}

	var accountingResult []model.DealerAccounting
	db.Raw("select dealers.id as dealer_id, order_lines.article_id as article_id, order_lines.price as price from dealers, orders, order_lines where orders.id = order_lines.order_id and dealers.id = ? ", vars["id"]).Scan(&accountingResult)

	if len(accountingResult) == 0 {
		w.WriteHeader(http.StatusNoContent)
		log.Printf("No PDFs to create")
		return
	}

	err = GenerateInvoicePdfHttp(w, accountingResult, dealerDetails)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("PDF Generator Error ", err)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
}
