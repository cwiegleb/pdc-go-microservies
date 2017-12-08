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

type dealerArticleCSV struct {
	Text     string  `csv:"Text"`
	Size     string  `csv:"Size"`
	Costs    float64 `csv:"Costs"`
	Currency string  `csv:"Currency"`
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

	multipartFileDealerArticles, _, err := r.FormFile("dealerArticles.csv")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print("failed to connect database", err)
		return
	}
	defer multipartFileDealerArticles.Close()

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

	dealerArticlesUpload := []*dealerArticleCSV{}
	if err := gocsv.Unmarshal(multipartFileDealerArticles, &dealerArticlesUpload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(err)
		return
	}

	dealer := &model.Dealer{
		Text:       dealerDetailsUpload[0].ExternalId,
		ExternalId: dealerDetailsUpload[0].ExternalId,
	}

	tx := db.Begin()
	defer tx.Close()

	var dealerWhere model.Dealer
	if db.Where("external_id = ?", dealer.ExternalId).First(&dealerWhere).Error == nil {
		tx.Rollback()
		w.WriteHeader(http.StatusBadRequest)
		log.Print("external id already exists")
		return
	}

	if tx.Create(dealer).Error != nil {
		tx.Rollback()
		w.WriteHeader(http.StatusUnprocessableEntity)
		log.Print("failed to create dealer")
	}

	dealerDetails := &model.DealerDetails{
		DealerID:   dealer.ID,
		Name:       dealerDetailsUpload[0].Name,
		Street:     dealerDetailsUpload[0].Street,
		City:       dealerDetailsUpload[0].City,
		PostalCode: dealerDetailsUpload[0].PostalCode,
		Telephone:  dealerDetailsUpload[0].Telephone,
		Email:      dealerDetailsUpload[0].Email,
		Iban:       dealerDetailsUpload[0].Iban,
		Bic:        dealerDetailsUpload[0].Bic,
		BankName:   dealerDetailsUpload[0].BankName,
		Fee:        dealerDetailsUpload[0].Fee,
		Commission: dealerDetailsUpload[0].Commission,
		Currency:   dealerDetailsUpload[0].Currency,
	}

	if tx.Create(dealerDetails).Error != nil {
		tx.Rollback()
		w.WriteHeader(http.StatusUnprocessableEntity)
		log.Print("failed to create dealer details")
	}

	for _, item := range dealerArticlesUpload {
		article := &model.Article{
			Text:      item.Text,
			Size:      item.Size,
			DealerID:  dealer.ID,
			Available: true,
			Costs:     item.Costs,
			Currency:  item.Currency,
		}

		if tx.Create(article).Error != nil {
			tx.Rollback()
			w.WriteHeader(http.StatusUnprocessableEntity)
			log.Print("failed to create article")
		}
	}
	tx.Commit()
}
