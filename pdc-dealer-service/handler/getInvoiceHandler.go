package handler

import (
	"log"
	"net/http"

	"github.com/cwiegleb/pdc-services/pdc-db/config"
	"github.com/cwiegleb/pdc-services/pdc-db/model"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/jung-kurt/gofpdf"
)

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

	var accountingResult []model.DealerAccounting
	db.Raw("select count(*) as number_sold_articles, sum(order_lines.price) as total_amount from dealers, orders, order_lines, articles where orders.id = order_lines.order_id and order_lines.article_id = articles.id and articles.dealer_id = dealers.id and dealers.id = ?", vars["id"]).Scan(&accountingResult)

	err = generatePDFInvoice(w, accountingResult)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("PDF Generator Error ", err)
	}

	w.Header().Set("Content-Type", "application/pdf")
}

func generatePDFInvoice(writer http.ResponseWriter, model []model.DealerAccounting) error {

	titleStr := "Rechnung"

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetTitle(titleStr, false)
	pdf.SetAuthor("PDC Kreischa", false)

	pdf.SetLeftMargin(15)
	pdf.SetRightMargin(15)
	pdf.SetTopMargin(20)
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 20)
	wd := pdf.GetStringWidth(titleStr) + 6
	pdf.CellFormat(wd, 14, titleStr, "0", 1, "L", false, 0, "")

	pdf.Ln(20)
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(0, 7, "Datum:", "0", 1, "R", false, 0, "")
	pdf.CellFormat(0, 7, "RECHNUNG NR:", "0", 1, "R", false, 0, "")

	pdf.Ln(10)
	pdf.CellFormat(0, 7, "An", "0", 1, "L", false, 0, "")
	pdf.CellFormat(0, 7, "[Vor- Nachname]", "0", 1, "L", false, 0, "")
	pdf.CellFormat(0, 7, "[Strasse / Hausnummer]", "0", 1, "L", false, 0, "")
	pdf.CellFormat(0, 7, "[PLZ / Ort]", "0", 1, "L", false, 0, "")
	pdf.CellFormat(0, 7, "[Telefon]", "0", 1, "L", false, 0, "")

	pdf.Ln(10)

	header := []string{"Position", "Artikel", "Summe der Position"}
	improvedTable := func() {
		// Column widths
		w := []float64{30.0, 90.0, 45.0}
		wSum := 0.0
		for _, v := range w {
			wSum += v
		}
		// 	Header
		for j, str := range header {
			pdf.CellFormat(w[j], 7, str, "1", 0, "C", false, 0, "")
		}
		pdf.Ln(-1)

		// Data
		//for _, c := range articlePositions {
		pdf.CellFormat(w[0], 7, "0", "LR", 0, "C", false, 0, "")
		pdf.CellFormat(w[1], 7, "Artikel 4711", "LR", 0, "", false, 0, "")
		pdf.CellFormat(w[2], 7, "47,11", "LR", 0, "C", false, 0, "")
		pdf.Ln(-1)
		//}
		pdf.SetFont("Arial", "B", 12)
		pdf.CellFormat(wSum, 0, "", "T", 1, "", false, 0, "")
		pdf.CellFormat(w[0]+w[1], 7, "Zwischensumme ", "", 0, "R", false, 0, "")
		pdf.CellFormat(w[2], 7, "47,11", "LR", 0, "C", false, 0, "")
		pdf.Ln(-1)
		pdf.CellFormat(w[0]+w[1], 0, "", "", 0, "", false, 0, "")
		pdf.CellFormat(w[2], 0, "", "T", 1, "", false, 0, "")

		pdf.CellFormat(w[0]+w[1], 7, "Provision ", "", 0, "R", false, 0, "")
		pdf.CellFormat(w[2], 7, "10%", "LR", 0, "C", false, 0, "")
		pdf.Ln(-1)
		pdf.CellFormat(w[0]+w[1], 0, "", "", 0, "", false, 0, "")
		pdf.CellFormat(w[2], 0, "", "T", 1, "", false, 0, "")
		pdf.SetFont("Arial", "", 12)
	}

	improvedTable()
	return pdf.Output(writer)
}
