package handler

import (
	"log"
	"net/http"
	"text/template"

	"github.com/cwiegleb/pdc-services/pdc-db/config"
	"github.com/jinzhu/gorm"
)

type DealerData struct {
	ExternalID string
	Iban       string
	Bic        string
	Name       string
	Amount     float32
}

type TransferAccounting struct {
	ExternalID string
	Price      float32
	Fee        float32
	Commission float32
	Iban       string
	Bic        string
	Name       string
}

func GetsHibiscusTransferHandler(w http.ResponseWriter, r *http.Request) {

	var config = config.LoadConfiguration("")
	db, err := gorm.Open(config.DbDriver, config.DbConnection)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print("failed to connect database", err)
		return
	}
	defer db.Close()

	var accountingResult []TransferAccounting
	db.Raw("select dealers.external_id, sum(order_lines.price) as price, dealer_details.fee, dealer_details.commission, dealer_details.iban, dealer_details.bic, dealer_details.name from dealers, orders, order_lines, dealer_details where orders.id = order_lines.order_id and order_lines.dealer_id = dealers.id and order_lines.dealer_id = dealer_details.dealer_id GROUP BY dealers.external_id, dealer_details.iban, dealer_details.bic, dealer_details.fee, dealer_details.commission, dealer_details.name   ORDER BY dealers.external_id").Scan(&accountingResult)
	if len(accountingResult) == 0 {
		w.WriteHeader(http.StatusNoContent)
		log.Printf("No PDFs to create")
		return
	}

	var dealerData = []DealerData{}
	for _, item := range accountingResult {
		dealerData = append(dealerData, DealerData{
			ExternalID: item.ExternalID,
			Iban:       item.Iban,
			Bic:        item.Bic,
			Amount:     item.Price - item.Fee - (item.Price * item.Commission / 100),
			Name:       item.Name,
		})
	}

	var template, tempErr = template.ParseFiles("./templates/hibiscus.transfer.tpl.xml")
	if err != nil {
		log.Print("failed to parse template", tempErr)
	}

	template.Execute(w, dealerData)

	w.Header().Set("Content-Type", "text/xml")
}
