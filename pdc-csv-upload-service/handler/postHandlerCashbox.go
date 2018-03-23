package handler

import (
	"log"
	"net/http"

	"github.com/cwiegleb/pdc-services/pdc-db/config"
	"github.com/cwiegleb/pdc-services/pdc-db/model"
	"github.com/gocarina/gocsv"
	"github.com/jinzhu/gorm"
)

type cashboxDataCSV struct {
	CashboxId  uint    `csv:"CashboxId"`
	ExternalId string  `csv:"ExternalId"`
	ArticleId  uint    `csv:"ArticleId"`
	Price      float64 `csv:"Price"`
	Currency   string  `csv:"Currency"`
}

func PostHandlerCashbox(w http.ResponseWriter, r *http.Request) {

	multipartFileCashboxData, _, err := r.FormFile("cashboxData.csv")
	defer multipartFileCashboxData.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print("Invalid input data", err)
		return
	}

	config := config.LoadConfiguration("")
	db, err := gorm.Open(config.DbDriver, config.DbConnection)
	defer db.Close()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print("failed to connect database", err)
		return
	}

	cashboxDataUpload := []*cashboxDataCSV{}
	if err := gocsv.Unmarshal(multipartFileCashboxData, &cashboxDataUpload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(err)
		return
	}

	tx := db.Begin()
	defer tx.Close()

	for _, item := range cashboxDataUpload {

		log.Print("unmarshal file", item)

		order := &model.Order{
			CashboxID:   item.CashboxId,
			OrderStatus: 1,
		}

		if tx.Create(order).Error != nil {
			tx.Rollback()
			w.WriteHeader(http.StatusUnprocessableEntity)
			log.Print("failed to create order")
			return
		}

		var dealerWhere model.Dealer
		if db.Where("external_id = ?", item.ExternalId).First(&dealerWhere).Error != nil {
			tx.Rollback()
			w.WriteHeader(http.StatusUnprocessableEntity)
			log.Print("failed to determine external id", item.ExternalId)
			return
		}

		var dealerDetailsWhere model.DealerDetails
		if db.Where("dealer_id = ?", dealerWhere.ID).First(&dealerDetailsWhere).Error == nil {
			tx.Rollback()
			w.WriteHeader(http.StatusUnprocessableEntity)
			log.Print("failed to determine dealer details for id", dealerDetailsWhere.ID)
			return
		}

		orderLine := &model.OrderLine{
			OrderID:    order.ID,
			ArticleID:  item.ArticleId,
			Price:      item.Price,
			DealerID:   dealerWhere.ID,
			DealerText: dealerWhere.ExternalId,
			Currency:   item.Currency,
		}

		if tx.Create(orderLine).Error != nil {
			tx.Rollback()
			w.WriteHeader(http.StatusUnprocessableEntity)
			log.Print("failed to create order line for id ", order.ID)
			return
		}
	}
	tx.Commit()
}
