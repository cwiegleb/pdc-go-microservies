package handler

import (
	"bufio"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cwiegleb/pdc-services/pdc-db/model"
	humanize "github.com/dustin/go-humanize"
	"github.com/jung-kurt/gofpdf"
)

/*
GenerateInvoicePdf Generate Pdf and write it to buffer
*/
func GenerateInvoicePdfHttp(writer http.ResponseWriter, dealerAccounting []model.DealerAccounting, dealerDetails model.DealerDetails) error {
	currentDateString := func() string {
		t := time.Now()
		return strings.Join([]string{strconv.Itoa(t.Day()), t.Month().String(), strconv.Itoa(t.Year())}, ".")
	}

	currentYearString := func() string {
		t := time.Now()
		return strconv.Itoa(t.Year())
	}

	titleStr := "Auszahlung"

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetTitle(titleStr, false)
	pdf.SetAuthor("PDC Kreischa", false)

	pdf.SetLeftMargin(15)
	pdf.SetRightMargin(15)
	pdf.SetTopMargin(20)
	pdf.AddPage()
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	pdf.SetFont("Helvetica", "B", 20)
	wd := pdf.GetStringWidth(titleStr) + 6
	pdf.CellFormat(wd, 14, titleStr, "0", 1, "L", false, 0, "")

	pdf.Ln(20)
	pdf.SetFont("Helvetica", "", 12)
	pdf.CellFormat(0, 7, currentDateString(), "0", 1, "R", false, 0, "")
	pdf.CellFormat(0, 7, strings.Join([]string{
		tr("Auszahlung: "), strings.Join([]string{
			currentYearString(),
			strconv.Itoa(int(time.Now().Month())),
			dealerAccounting[0].ExternalID}, "_")}, " "), "0", 1, "R", false, 0, "")

	pdf.Ln(10)
	pdf.CellFormat(0, 7, tr("An"), "0", 1, "L", false, 0, "")
	pdf.CellFormat(0, 7, tr(dealerDetails.Name), "0", 1, "L", false, 0, "")
	if len(dealerDetails.Street) != 0 && len(dealerDetails.City) != 0 {
		pdf.CellFormat(0, 7, tr(dealerDetails.Street), "0", 1, "L", false, 0, "")
		pdf.CellFormat(0, 7, strings.Join([]string{tr(dealerDetails.PostalCode), tr(dealerDetails.City)}, " "), "0", 1, "L", false, 0, "")
	}
	pdf.Ln(10)
	pdf.CellFormat(0, 7, strings.Join([]string{"Anbieternummer", tr(dealerAccounting[0].ExternalID)}, " "), "0", 1, "L", false, 0, "")

	pdf.Ln(10)

	header := []string{"Position", "Artikel", "Summe der Position"}

	partAmount := func() float32 {
		var amount float32
		for _, item := range dealerAccounting {
			amount = amount + item.Price
		}
		return amount
	}

	partAmountCommission := func() float32 {
		return partAmount() * dealerDetails.Commission / 100
	}

	toalAmount := func() float32 {
		return partAmount() - partAmountCommission() - dealerDetails.Fee
	}

	improvedTable := func() {
		// Column widths
		w := []float64{30.0, 100.0, 45.0}
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
		for j, item := range dealerAccounting {
			pdf.CellFormat(w[0], 7, strconv.Itoa(int(j)+1), "LR", 0, "C", false, 0, "")
			pdf.CellFormat(w[1], 7, strconv.Itoa(int(item.ArticleID)), "LR", 0, "C", false, 0, "")
			pdf.CellFormat(w[2], 7, humanize.FormatFloat("###,##", float64(item.Price)), "LR", 0, "C", false, 0, "")
			pdf.Ln(-1)
		}

		pdf.SetFont("Helvetica", "B", 12)
		pdf.CellFormat(wSum, 0, "", "T", 1, "", false, 0, "")
		pdf.CellFormat(w[0]+w[1], 7, tr("Zwischensumme in € "), "", 0, "R", false, 0, "")
		pdf.CellFormat(w[2], 7, humanize.FormatFloat("###,##", float64(partAmount())), "LR", 0, "C", false, 0, "")
		pdf.Ln(-1)
		pdf.CellFormat(w[0]+w[1], 0, "", "", 0, "", false, 0, "")
		pdf.CellFormat(w[2], 0, "", "T", 1, "", false, 0, "")

		pdf.CellFormat(w[0]+w[1], 7, tr("Abzüglich Unterstützung Förderverein in % "), "", 0, "R", false, 0, "")
		pdf.CellFormat(w[2], 7, humanize.FormatFloat("###,##", float64(dealerDetails.Commission)), "LR", 0, "C", false, 0, "")
		pdf.Ln(-1)
		pdf.CellFormat(w[0]+w[1], 0, "", "", 0, "", false, 0, "")
		pdf.CellFormat(w[2], 0, "", "T", 1, "", false, 0, "")

		pdf.CellFormat(w[0]+w[1], 7, tr("Abzüglich Unterstützung Förderverein in € "), "", 0, "R", false, 0, "")
		pdf.CellFormat(w[2], 7, humanize.FormatFloat("###,##", float64(partAmountCommission())), "LR", 0, "C", false, 0, "")
		pdf.Ln(-1)
		pdf.CellFormat(w[0]+w[1], 0, "", "", 0, "", false, 0, "")
		pdf.CellFormat(w[2], 0, "", "T", 1, "", false, 0, "")

		pdf.CellFormat(w[0]+w[1], 7, tr("Abzüglich Pauschale für Abwicklung in € "), "", 0, "R", false, 0, "")
		pdf.CellFormat(w[2], 7, humanize.FormatFloat("###,##", float64(dealerDetails.Fee)), "LR", 0, "C", false, 0, "")
		pdf.Ln(-1)
		pdf.CellFormat(w[0]+w[1], 0, "", "", 0, "", false, 0, "")
		pdf.CellFormat(w[2], 0, "", "T", 1, "", false, 0, "")

		pdf.CellFormat(w[0]+w[1], 7, tr("Auszahlung in € "), "", 0, "R", false, 0, "")
		pdf.CellFormat(w[2], 7, humanize.FormatFloat("###,##", float64(toalAmount())), "LR", 0, "C", false, 0, "")
		pdf.Ln(-1)
		pdf.CellFormat(w[0]+w[1], 0, "", "", 0, "", false, 0, "")
		pdf.CellFormat(w[2], 0, "", "T", 1, "", false, 0, "")

		pdf.SetFont("Helvetica", "", 12)

		// if len(dealerDetails.Iban) != 0 {
		// 	pdf.Ln(20)
		// 	pdf.CellFormat(0, 7, tr("Die Auszahlung erfolgt auf folgendes Konto:"), "0", 1, "L", false, 0, "")
		// 	pdf.CellFormat(0, 7, strings.Join([]string{dealerDetails.Iban, dealerDetails.Bic}, " "), "0", 1, "L", false, 0, "")
		// }

		pdf.Ln(20)
		pdf.CellFormat(0, 7, tr("Mit freundlichen Grüßen"), "0", 1, "L", false, 0, "")
		pdf.CellFormat(0, 7, tr("Kinderkleiderbörse Kreischa"), "0", 1, "L", false, 0, "")
	}

	improvedTable()

	return pdf.Output(writer)
}

/*
GenerateInvoicePdf Generate Pdf and write it to buffer
*/
func GenerateInvoicePdfBuffer(writer *bufio.Writer, dealerAccounting []model.DealerAccounting, dealerDetails model.DealerDetails) error {
	currentDateString := func() string {
		t := time.Now()
		return strings.Join([]string{strconv.Itoa(t.Day()), t.Month().String(), strconv.Itoa(t.Year())}, ".")
	}

	currentYearString := func() string {
		t := time.Now()
		return strconv.Itoa(t.Year())
	}

	titleStr := "Auszahlung"

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetTitle(titleStr, false)
	pdf.SetAuthor("PDC Kreischa", false)

	pdf.SetLeftMargin(15)
	pdf.SetRightMargin(15)
	pdf.SetTopMargin(20)
	pdf.AddPage()
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	pdf.SetFont("Helvetica", "B", 20)
	wd := pdf.GetStringWidth(titleStr) + 6
	pdf.CellFormat(wd, 14, titleStr, "0", 1, "L", false, 0, "")

	pdf.Ln(20)
	pdf.SetFont("Helvetica", "", 12)
	pdf.CellFormat(0, 7, currentDateString(), "0", 1, "R", false, 0, "")
	pdf.CellFormat(0, 7, strings.Join([]string{
		tr("Auszahlung: "), strings.Join([]string{
			currentYearString(),
			strconv.Itoa(int(time.Now().Month())),
			dealerAccounting[0].ExternalID}, "_")}, " "), "0", 1, "R", false, 0, "")

	pdf.Ln(10)
	pdf.CellFormat(0, 7, tr("An"), "0", 1, "L", false, 0, "")
	pdf.CellFormat(0, 7, tr(dealerDetails.Name), "0", 1, "L", false, 0, "")
	if len(dealerDetails.Street) != 0 && len(dealerDetails.City) != 0 {
		pdf.CellFormat(0, 7, tr(dealerDetails.Street), "0", 1, "L", false, 0, "")
		pdf.CellFormat(0, 7, strings.Join([]string{tr(dealerDetails.PostalCode), tr(dealerDetails.City)}, " "), "0", 1, "L", false, 0, "")
	}
	pdf.Ln(10)
	pdf.CellFormat(0, 7, strings.Join([]string{"Anbieternummer", tr(dealerAccounting[0].ExternalID)}, " "), "0", 1, "L", false, 0, "")

	pdf.Ln(10)

	header := []string{"Position", "Artikel", "Summe der Position"}

	partAmount := func() float32 {
		var amount float32
		for _, item := range dealerAccounting {
			amount = amount + item.Price
		}
		return amount
	}

	partAmountCommission := func() float32 {
		return partAmount() * dealerDetails.Commission / 100
	}

	toalAmount := func() float32 {
		return partAmount() - partAmountCommission() - dealerDetails.Fee
	}

	improvedTable := func() {
		// Column widths
		w := []float64{30.0, 100.0, 45.0}
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
		for j, item := range dealerAccounting {
			pdf.CellFormat(w[0], 7, strconv.Itoa(int(j)+1), "LR", 0, "C", false, 0, "")
			pdf.CellFormat(w[1], 7, strconv.Itoa(int(item.ArticleID)), "LR", 0, "C", false, 0, "")
			pdf.CellFormat(w[2], 7, humanize.FormatFloat("###,##", float64(item.Price)), "LR", 0, "C", false, 0, "")
			pdf.Ln(-1)
		}

		pdf.SetFont("Helvetica", "B", 12)
		pdf.CellFormat(wSum, 0, "", "T", 1, "", false, 0, "")
		pdf.CellFormat(w[0]+w[1], 7, tr("Zwischensumme in € "), "", 0, "R", false, 0, "")
		pdf.CellFormat(w[2], 7, humanize.FormatFloat("###,##", float64(partAmount())), "LR", 0, "C", false, 0, "")
		pdf.Ln(-1)
		pdf.CellFormat(w[0]+w[1], 0, "", "", 0, "", false, 0, "")
		pdf.CellFormat(w[2], 0, "", "T", 1, "", false, 0, "")

		pdf.CellFormat(w[0]+w[1], 7, tr("Abzüglich Unterstützung Förderverein in % "), "", 0, "R", false, 0, "")
		pdf.CellFormat(w[2], 7, humanize.FormatFloat("###,##", float64(dealerDetails.Commission)), "LR", 0, "C", false, 0, "")
		pdf.Ln(-1)
		pdf.CellFormat(w[0]+w[1], 0, "", "", 0, "", false, 0, "")
		pdf.CellFormat(w[2], 0, "", "T", 1, "", false, 0, "")

		pdf.CellFormat(w[0]+w[1], 7, tr("Abzüglich Unterstützung Förderverein in € "), "", 0, "R", false, 0, "")
		pdf.CellFormat(w[2], 7, humanize.FormatFloat("###,##", float64(partAmountCommission())), "LR", 0, "C", false, 0, "")
		pdf.Ln(-1)
		pdf.CellFormat(w[0]+w[1], 0, "", "", 0, "", false, 0, "")
		pdf.CellFormat(w[2], 0, "", "T", 1, "", false, 0, "")

		pdf.CellFormat(w[0]+w[1], 7, tr("Abzüglich Pauschale für Abwicklung in € "), "", 0, "R", false, 0, "")
		pdf.CellFormat(w[2], 7, humanize.FormatFloat("###,##", float64(dealerDetails.Fee)), "LR", 0, "C", false, 0, "")
		pdf.Ln(-1)
		pdf.CellFormat(w[0]+w[1], 0, "", "", 0, "", false, 0, "")
		pdf.CellFormat(w[2], 0, "", "T", 1, "", false, 0, "")

		pdf.CellFormat(w[0]+w[1], 7, tr("Auszahlung in € "), "", 0, "R", false, 0, "")
		pdf.CellFormat(w[2], 7, humanize.FormatFloat("###,##", float64(toalAmount())), "LR", 0, "C", false, 0, "")
		pdf.Ln(-1)
		pdf.CellFormat(w[0]+w[1], 0, "", "", 0, "", false, 0, "")
		pdf.CellFormat(w[2], 0, "", "T", 1, "", false, 0, "")

		pdf.SetFont("Helvetica", "", 12)

		// if len(dealerDetails.Iban) != 0 {
		// 	pdf.Ln(20)
		// 	pdf.CellFormat(0, 7, tr("Die Auszahlung erfolgt auf folgendes Konto:"), "0", 1, "L", false, 0, "")
		// 	pdf.CellFormat(0, 7, strings.Join([]string{dealerDetails.Iban, dealerDetails.Bic}, " "), "0", 1, "L", false, 0, "")
		// }

		pdf.Ln(20)
		pdf.CellFormat(0, 7, tr("Mit freundlichen Grüßen"), "0", 1, "L", false, 0, "")
		pdf.CellFormat(0, 7, tr("Kinderkleiderbörse Kreischa"), "0", 1, "L", false, 0, "")
	}

	improvedTable()
	return pdf.Output(writer)
}
