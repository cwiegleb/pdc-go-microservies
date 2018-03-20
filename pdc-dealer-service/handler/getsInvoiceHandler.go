package handler

import (
	"archive/zip"
	"bufio"
	"bytes"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/cwiegleb/pdc-services/pdc-db/config"
	"github.com/cwiegleb/pdc-services/pdc-db/model"
	"github.com/jinzhu/gorm"
)

/*
GetsInvoiceHandler
*/
func GetsInvoiceHandler(w http.ResponseWriter, r *http.Request) {

	var config = config.LoadConfiguration("")
	db, err := gorm.Open(config.DbDriver, config.DbConnection)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print("Failed to connect database", err)
		return
	}
	defer db.Close()

	var dealersDetails []model.DealerDetails
	if db.Find(&dealersDetails).Error != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Print("No entries found")
		return
	}

	var accountingResult []model.DealerAccounting
	db.Raw("select dealers.id as dealer_id, order_lines.article_id as article_id, order_lines.price as price from dealers, orders, order_lines where orders.id = order_lines.order_id and order_lines.dealer_id = dealers.id").Scan(&accountingResult)
	if len(accountingResult) == 0 {
		w.WriteHeader(http.StatusNoContent)
		log.Printf("No PDFs to create")
		return
	}
	accountingResultMap := make(map[int][]model.DealerAccounting)

	for _, item := range accountingResult {
		accountingResultMap[int(item.DealerID)] = append(accountingResultMap[int(item.DealerID)], item)
	}

	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	var fw io.Writer

	for _, item := range accountingResultMap {
		var pdfBytes bytes.Buffer
		writer := bufio.NewWriter(&pdfBytes)
		var dealerDetail model.DealerDetails

		for _, dealerItem := range dealersDetails {
			if dealerItem.DealerID == item[0].DealerID {
				dealerDetail = dealerItem
				break
			}
		}

		err = GenerateInvoicePdfBuffer(writer, item, dealerDetail)
		writer.Flush()

		if fw, err = zipWriter.Create("Auszahlung_" + strconv.Itoa(int(item[0].DealerID)) + ".pdf"); err != nil {
			log.Fatal(err)
		}
		if _, err = fw.Write(pdfBytes.Bytes()); err != nil {
			log.Fatal(err)
		}
	}
	w.Header().Set("Content-Type", "application/zip")
}
