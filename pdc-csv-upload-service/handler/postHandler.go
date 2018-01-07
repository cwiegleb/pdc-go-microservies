package handler

import (
	"log"
	"net/http"

	"github.com/cwiegleb/pdc-services/pdc-db/config"
	"github.com/cwiegleb/pdc-services/pdc-db/model"
	"github.com/gocarina/gocsv"
	"github.com/jinzhu/gorm"
)

type dealerDetailsCSV struct {
	ExternalId string  `csv:"ExternalId"`
	Text       string  `csv:"Text"`
	Name       string  `csv:"Name"`
	Street     string  `csv:"Street"`
	City       string  `csv:"City"`
	PostalCode string  `csv:"PostalCode"`
	Telephone  string  `csv:"Telephone"`
	Email      string  `csv:"Email"`
	Iban       string  `csv:"Iban"`
	Bic        string  `csv:"Bic"`
	BankName   string  `csv:"BankName"`
	Fee        float32 `csv:"Fee"`
	Commission float32 `csv:"Commission"`
	Currency   string  `csv:"Currency"`
}

type errorMessage struct {
	message string
}

func PostHandler(w http.ResponseWriter, r *http.Request) {

	multipartFileDealer, _, err := r.FormFile("dealerDetails.csv")
	defer multipartFileDealer.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print("failed to connect database", err)
		return
	}
	defer multipartFileDealer.Close()

	config := config.LoadConfiguration("")
	db, err := gorm.Open(config.DbDriver, config.DbConnection)
	defer db.Close()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print("failed to connect database", err)
		return
	}

	dealerDetailsUpload := []*dealerDetailsCSV{}
	if err := gocsv.Unmarshal(multipartFileDealer, &dealerDetailsUpload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(err)
		return
	}

	tx := db.Begin()
	defer tx.Close()

	for _, item := range dealerDetailsUpload {

		dealer := &model.Dealer{
			Text:       item.ExternalId,
			ExternalId: item.ExternalId,
		}

		var dealerWhere model.Dealer
		if db.Where("external_id = ?", dealer.ExternalId).First(&dealerWhere).Error == nil {

			var dealerDetailsWhere model.DealerDetails
			if db.Where("dealer_id = ?", dealerWhere.ID).First(&dealerDetailsWhere).Error == nil {

				dealerDetailsWhere.DealerID = dealerWhere.ID
				dealerDetailsWhere.Name = item.Name
				dealerDetailsWhere.Street = item.Street
				dealerDetailsWhere.City = item.City
				dealerDetailsWhere.PostalCode = item.PostalCode
				dealerDetailsWhere.Telephone = item.Telephone
				dealerDetailsWhere.Email = item.Email
				dealerDetailsWhere.Iban = item.Iban
				dealerDetailsWhere.Bic = item.Bic
				dealerDetailsWhere.BankName = item.BankName
				dealerDetailsWhere.Fee = item.Fee
				dealerDetailsWhere.Commission = item.Commission
				dealerDetailsWhere.Currency = item.Currency

				if tx.Save(dealerDetailsWhere).Error != nil {
					tx.Rollback()
					w.WriteHeader(http.StatusUnprocessableEntity)
					log.Print("failed to update dealer with id", dealerWhere.ID)
					return
				}
			}

		} else {
			if tx.Create(dealer).Error != nil {
				tx.Rollback()
				w.WriteHeader(http.StatusUnprocessableEntity)
				log.Print("failed to create dealer")
				return
			}

			dealerDetails := &model.DealerDetails{
				DealerID:   dealer.ID,
				Name:       item.Name,
				Street:     item.Street,
				City:       item.City,
				PostalCode: item.PostalCode,
				Telephone:  item.Telephone,
				Email:      item.Email,
				Iban:       item.Iban,
				Bic:        item.Bic,
				BankName:   item.BankName,
				Fee:        item.Fee,
				Commission: item.Commission,
				Currency:   item.Currency,
			}

			if tx.Create(dealerDetails).Error != nil {
				tx.Rollback()
				w.WriteHeader(http.StatusUnprocessableEntity)
				log.Print("failed to create dealer details for id ", dealerDetails.DealerID)
				return
			}
		}
	}
	tx.Commit()
}
