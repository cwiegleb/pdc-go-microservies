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
	"github.com/jung-kurt/gofpdf"
)

func GetsInvoiceHandler(w http.ResponseWriter, r *http.Request) {

	var config = config.LoadConfiguration("")
	db, err := gorm.Open(config.DbDriver, config.DbConnection)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print("failed to connect database", err)
		return
	}
	defer db.Close()

	var accountingResult []model.DealerAccounting
	db.Raw("select dealers.id as dealer_id, count(*) as number_sold_articles, sum(order_lines.price) as total_amount from dealers, orders, order_lines, articles where orders.id = order_lines.order_id and order_lines.article_id = articles.id and articles.dealer_id = dealers.id group by dealers.id").Scan(&accountingResult)

	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	var fw io.Writer

	for _, item := range accountingResult {
		var pdfBytes bytes.Buffer
		writer := bufio.NewWriter(&pdfBytes)
		err = generateBufferPDFInvoice(writer, item)
		writer.Flush()

		if fw, err = zipWriter.Create("Test_" + strconv.Itoa(int(item.DealerID)) + ".pdf"); err != nil {
			log.Fatal(err)
		}
		if _, err = fw.Write(pdfBytes.Bytes()); err != nil {
			log.Fatal(err)
		}
	}
	w.Header().Set("Content-Type", "application/zip")
}

func generateBufferPDFInvoice(writer *bufio.Writer, model model.DealerAccounting) error {

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetLeftMargin(45)
	pdf.SetFont("Arial", "B", 14)
	_, lineHt := pdf.GetFontSize()
	htmlStr :=
		`
		<style>
		<!--
		/* Font Definitions */
		@font-face
			{font-family:Arial;
			panose-1:2 11 6 4 2 2 2 2 2 4;
			mso-font-charset:0;
			mso-generic-font-family:auto;
			mso-font-pitch:variable;
			mso-font-signature:-536859905 -1073711037 9 0 511 0;}
		@font-face
			{font-family:"Cambria Math";
			panose-1:2 4 5 3 5 4 6 3 2 4;
			mso-font-charset:1;
			mso-generic-font-family:roman;
			mso-font-format:other;
			mso-font-pitch:variable;
			mso-font-signature:0 0 0 0 0 0;}
		@font-face
			{font-family:SimSun;
			panose-1:2 1 6 0 3 1 1 1 1 1;
			mso-font-alt:__;
			mso-font-charset:134;
			mso-generic-font-family:auto;
			mso-font-format:other;
			mso-font-pitch:variable;
			mso-font-signature:1 135135232 16 0 262144 0;}
		@font-face
			{font-family:"Microsoft Sans Serif";
			panose-1:2 11 6 4 2 2 2 2 2 4;
			mso-font-charset:0;
			mso-generic-font-family:auto;
			mso-font-pitch:variable;
			mso-font-signature:-520082689 -1073741822 8 0 66047 0;}
		@font-face
			{font-family:Tahoma;
			panose-1:2 11 6 4 3 5 4 4 2 4;
			mso-font-charset:0;
			mso-generic-font-family:swiss;
			mso-font-pitch:variable;
			mso-font-signature:1627400839 -2147483648 8 0 66047 0;}
		/* Style Definitions */
		p.MsoNormal, li.MsoNormal, div.MsoNormal
			{mso-style-unhide:no;
			mso-style-qformat:yes;
			mso-style-parent:"";
			margin:0cm;
			margin-bottom:.0001pt;
			mso-pagination:widow-orphan;
			font-size:8.0pt;
			mso-bidi-font-size:12.0pt;
			font-family:"Microsoft Sans Serif";
			mso-ascii-font-family:"Microsoft Sans Serif";
			mso-ascii-theme-font:minor-latin;
			mso-fareast-font-family:SimSun;
			mso-hansi-font-family:"Microsoft Sans Serif";
			mso-hansi-theme-font:minor-latin;
			mso-bidi-font-family:"Times New Roman";
			mso-ansi-language:EN-US;
			mso-fareast-language:EN-US;}
		h1
			{mso-style-update:auto;
			mso-style-unhide:no;
			mso-style-qformat:yes;
			mso-style-next:Standard;
			margin:0cm;
			margin-bottom:.0001pt;
			text-align:right;
			line-height:40.0pt;
			mso-line-height-rule:exactly;
			mso-pagination:widow-orphan;
			page-break-after:avoid;
			mso-outline-level:1;
			font-size:36.0pt;
			font-family:"Microsoft Sans Serif";
			mso-ascii-font-family:"Microsoft Sans Serif";
			mso-ascii-theme-font:major-latin;
			mso-hansi-font-family:"Microsoft Sans Serif";
			mso-hansi-theme-font:major-latin;
			mso-bidi-font-family:Arial;
			color:#DFEADF;
			mso-themecolor:accent2;
			mso-themetint:102;
			mso-font-kerning:22.0pt;
			mso-ansi-language:EN-US;
			mso-fareast-language:EN-US;
			font-weight:normal;
			mso-bidi-font-weight:bold;}
		h2
			{mso-style-unhide:no;
			mso-style-qformat:yes;
			mso-style-next:Standard;
			margin-top:12.0pt;
			margin-right:0cm;
			margin-bottom:3.0pt;
			margin-left:0cm;
			mso-pagination:widow-orphan;
			page-break-after:avoid;
			mso-outline-level:2;
			font-size:14.0pt;
			font-family:"Microsoft Sans Serif";
			mso-ascii-font-family:"Microsoft Sans Serif";
			mso-ascii-theme-font:major-latin;
			mso-hansi-font-family:"Microsoft Sans Serif";
			mso-hansi-theme-font:major-latin;
			mso-bidi-font-family:Arial;
			mso-ansi-language:EN-US;
			mso-fareast-language:EN-US;
			font-style:italic;}
		h3
			{mso-style-unhide:no;
			mso-style-qformat:yes;
			mso-style-next:Standard;
			margin-top:12.0pt;
			margin-right:0cm;
			margin-bottom:3.0pt;
			margin-left:0cm;
			mso-pagination:widow-orphan;
			page-break-after:avoid;
			mso-outline-level:3;
			font-size:13.0pt;
			font-family:"Microsoft Sans Serif";
			mso-ascii-font-family:"Microsoft Sans Serif";
			mso-ascii-theme-font:major-latin;
			mso-hansi-font-family:"Microsoft Sans Serif";
			mso-hansi-theme-font:major-latin;
			mso-bidi-font-family:Arial;
			mso-ansi-language:EN-US;
			mso-fareast-language:EN-US;}
		p.MsoAcetate, li.MsoAcetate, div.MsoAcetate
			{mso-style-unhide:no;
			mso-style-link:"Sprechblasentext Zchn";
			margin:0cm;
			margin-bottom:.0001pt;
			mso-pagination:widow-orphan;
			font-size:8.0pt;
			font-family:"Tahoma","sans-serif";
			mso-fareast-font-family:SimSun;
			mso-ansi-language:EN-US;
			mso-fareast-language:EN-US;}
		span.MsoPlaceholderText
			{mso-style-noshow:yes;
			mso-style-priority:99;
			mso-style-unhide:no;
			color:gray;}
		p.Betrag, li.Betrag, div.Betrag
			{mso-style-name:Betrag;
			mso-style-unhide:no;
			margin:0cm;
			margin-bottom:.0001pt;
			text-align:right;
			mso-pagination:widow-orphan;
			font-size:8.0pt;
			mso-bidi-font-size:12.0pt;
			font-family:"Microsoft Sans Serif";
			mso-ascii-font-family:"Microsoft Sans Serif";
			mso-ascii-theme-font:minor-latin;
			mso-fareast-font-family:SimSun;
			mso-hansi-font-family:"Microsoft Sans Serif";
			mso-hansi-theme-font:minor-latin;
			mso-bidi-font-family:"Times New Roman";
			mso-ansi-language:EN-US;
			mso-fareast-language:EN-US;}
		p.DatumundNummer, li.DatumundNummer, div.DatumundNummer
			{mso-style-name:"Datum und Nummer";
			mso-style-unhide:no;
			margin:0cm;
			margin-bottom:.0001pt;
			text-align:right;
			line-height:110%;
			mso-pagination:widow-orphan;
			font-size:8.0pt;
			font-family:"Microsoft Sans Serif";
			mso-ascii-font-family:"Microsoft Sans Serif";
			mso-ascii-theme-font:minor-latin;
			mso-fareast-font-family:SimSun;
			mso-hansi-font-family:"Microsoft Sans Serif";
			mso-hansi-theme-font:minor-latin;
			mso-bidi-font-family:"Times New Roman";
			color:gray;
			mso-themecolor:background1;
			mso-themeshade:128;
			letter-spacing:.2pt;
			mso-ansi-language:EN-US;
			mso-fareast-language:EN-US;
			font-weight:bold;
			mso-bidi-font-weight:normal;}
		p.berschriften, li.berschriften, div.berschriften
			{mso-style-name:�berschriften;
			mso-style-unhide:no;
			mso-style-parent:"Rechts ausgerichteter Text";
			margin:0cm;
			margin-bottom:.0001pt;
			text-align:right;
			line-height:12.0pt;
			mso-pagination:widow-orphan;
			font-size:8.0pt;
			font-family:"Microsoft Sans Serif";
			mso-ascii-font-family:"Microsoft Sans Serif";
			mso-ascii-theme-font:major-latin;
			mso-fareast-font-family:SimSun;
			mso-hansi-font-family:"Microsoft Sans Serif";
			mso-hansi-theme-font:major-latin;
			mso-bidi-font-family:"Times New Roman";
			color:gray;
			mso-themecolor:background1;
			mso-themeshade:128;
			text-transform:uppercase;
			letter-spacing:.2pt;
			mso-ansi-language:EN-US;
			mso-fareast-language:EN-US;
			font-weight:bold;}
		p.Slogan, li.Slogan, div.Slogan
			{mso-style-name:Slogan;
			mso-style-unhide:no;
			mso-style-parent:"�berschrift 3";
			margin-top:0cm;
			margin-right:0cm;
			margin-bottom:3.0pt;
			margin-left:0cm;
			mso-pagination:widow-orphan;
			mso-outline-level:3;
			font-size:8.0pt;
			mso-bidi-font-size:9.0pt;
			font-family:"Microsoft Sans Serif";
			mso-ascii-font-family:"Microsoft Sans Serif";
			mso-ascii-theme-font:minor-latin;
			mso-fareast-font-family:SimSun;
			mso-hansi-font-family:"Microsoft Sans Serif";
			mso-hansi-theme-font:minor-latin;
			mso-bidi-font-family:"Times New Roman";
			color:gray;
			mso-themecolor:background1;
			mso-themeshade:128;
			letter-spacing:.2pt;
			mso-ansi-language:EN-US;
			mso-fareast-language:EN-US;
			font-weight:bold;
			mso-bidi-font-weight:normal;
			font-style:italic;
			mso-bidi-font-style:normal;}
		p.ZentrierterTextunten, li.ZentrierterTextunten, div.ZentrierterTextunten
			{mso-style-name:"Zentrierter Text unten";
			mso-style-unhide:no;
			margin-top:26.0pt;
			margin-right:0cm;
			margin-bottom:0cm;
			margin-left:0cm;
			margin-bottom:.0001pt;
			text-align:center;
			mso-pagination:widow-orphan;
			font-size:8.0pt;
			mso-bidi-font-size:9.0pt;
			font-family:"Microsoft Sans Serif";
			mso-ascii-font-family:"Microsoft Sans Serif";
			mso-ascii-theme-font:minor-latin;
			mso-fareast-font-family:SimSun;
			mso-hansi-font-family:"Microsoft Sans Serif";
			mso-hansi-theme-font:minor-latin;
			mso-bidi-font-family:"Times New Roman";
			color:#B0CCB0;
			mso-themecolor:accent2;
			mso-ansi-language:EN-US;
			mso-fareast-language:EN-US;}
		p.Spaltenberschrift, li.Spaltenberschrift, div.Spaltenberschrift
			{mso-style-name:Spalten�berschrift;
			mso-style-update:auto;
			mso-style-unhide:no;
			mso-style-parent:"�berschrift 2";
			margin-top:1.0pt;
			margin-right:0cm;
			margin-bottom:0cm;
			margin-left:0cm;
			margin-bottom:.0001pt;
			mso-pagination:widow-orphan;
			mso-outline-level:2;
			font-size:8.0pt;
			font-family:"Microsoft Sans Serif";
			mso-ascii-font-family:"Microsoft Sans Serif";
			mso-ascii-theme-font:major-latin;
			mso-fareast-font-family:SimSun;
			mso-hansi-font-family:"Microsoft Sans Serif";
			mso-hansi-theme-font:major-latin;
			mso-bidi-font-family:"Times New Roman";
			color:gray;
			mso-themecolor:background1;
			mso-themeshade:128;
			mso-ansi-language:EN-US;
			mso-fareast-language:EN-US;
			font-weight:bold;
			mso-bidi-font-weight:normal;}
		span.SprechblasentextZchn
			{mso-style-name:"Sprechblasentext Zchn";
			mso-style-unhide:no;
			mso-style-locked:yes;
			mso-style-link:Sprechblasentext;
			mso-ansi-font-size:8.0pt;
			mso-bidi-font-size:8.0pt;
			font-family:"Tahoma","sans-serif";
			mso-ascii-font-family:Tahoma;
			mso-hansi-font-family:Tahoma;
			mso-bidi-font-family:Tahoma;}
		p.VielenDank, li.VielenDank, div.VielenDank
			{mso-style-name:"Vielen Dank";
			mso-style-update:auto;
			mso-style-unhide:no;
			margin-top:5.0pt;
			margin-right:0cm;
			margin-bottom:0cm;
			margin-left:0cm;
			margin-bottom:.0001pt;
			text-align:center;
			mso-pagination:widow-orphan;
			font-size:10.0pt;
			mso-bidi-font-size:12.0pt;
			font-family:"Microsoft Sans Serif";
			mso-ascii-font-family:"Microsoft Sans Serif";
			mso-ascii-theme-font:minor-latin;
			mso-fareast-font-family:SimSun;
			mso-hansi-font-family:"Microsoft Sans Serif";
			mso-hansi-theme-font:minor-latin;
			mso-bidi-font-family:"Times New Roman";
			color:gray;
			mso-themecolor:background1;
			mso-themeshade:128;
			mso-ansi-language:EN-US;
			mso-fareast-language:EN-US;
			font-style:italic;
			mso-bidi-font-style:normal;}
		p.RechtsausgerichteterText, li.RechtsausgerichteterText, div.RechtsausgerichteterText
			{mso-style-name:"Rechts ausgerichteter Text";
			mso-style-unhide:no;
			margin:0cm;
			margin-bottom:.0001pt;
			text-align:right;
			line-height:12.0pt;
			mso-pagination:widow-orphan;
			font-size:8.0pt;
			font-family:"Microsoft Sans Serif";
			mso-ascii-font-family:"Microsoft Sans Serif";
			mso-ascii-theme-font:minor-latin;
			mso-fareast-font-family:SimSun;
			mso-hansi-font-family:"Microsoft Sans Serif";
			mso-hansi-theme-font:minor-latin;
			mso-bidi-font-family:"Times New Roman";
			color:#7F909D;
			mso-ansi-language:EN-US;
			mso-fareast-language:EN-US;
			font-weight:bold;
			mso-bidi-font-weight:normal;}
		.MsoChpDefault
			{mso-style-type:export-only;
			mso-default-props:yes;
			font-size:10.0pt;
			mso-ansi-font-size:10.0pt;
			mso-bidi-font-size:10.0pt;
			mso-fareast-font-family:SimSun;
			mso-ansi-language:EN-US;
			mso-fareast-language:EN-US;}
		@page WordSection1
			{size:21.0cm 841.95pt;
			margin:72.0pt 72.0pt 72.0pt 72.0pt;
			mso-header-margin:36.0pt;
			mso-footer-margin:36.0pt;
			mso-paper-source:0;}
		div.WordSection1
			{page:WordSection1;}
		-->
		</style>
		<!--[if gte mso 10]>
		<style>
		/* Style Definitions */
		table.MsoNormalTable
			{mso-style-name:"Normale Tabelle";
			mso-tstyle-rowband-size:0;
			mso-tstyle-colband-size:0;
			mso-style-noshow:yes;
			mso-style-priority:99;
			mso-style-parent:"";
			mso-padding-alt:0cm 5.4pt 0cm 5.4pt;
			mso-para-margin:0cm;
			mso-para-margin-bottom:.0001pt;
			mso-pagination:widow-orphan;
			font-size:10.0pt;
			font-family:"Times New Roman";
			mso-ansi-language:EN-US;
			mso-fareast-language:EN-US;}
		</style>
		<![endif]--><!--[if gte mso 9]><xml>
		<o:shapedefaults v:ext="edit" spidmax="1026" fillcolor="white">
		<v:fill color="white"/>
		</o:shapedefaults></xml><![endif]--><!--[if gte mso 9]><xml>
		<o:shapelayout v:ext="edit">
		<o:idmap v:ext="edit" data="1"/>
		</o:shapelayout></xml><![endif]-->
		</head>
		
		<body lang=DE style='tab-interval:36.0pt'>
		
		<div class=WordSection1>
		
		<div align=center>
		
		<table class=MsoNormalTable border=0 cellspacing=0 cellpadding=0 width=504
		style='width:504.0pt;border-collapse:collapse;mso-padding-alt:7.2pt 5.75pt 2.9pt 5.75pt'>
		<tr style='mso-yfti-irow:0;mso-yfti-firstrow:yes;height:37.8pt'>
		<td width=79 valign=top style='width:79.15pt;padding:0cm 5.75pt 2.9pt 5.75pt;
		height:37.8pt'>
		<p class=MsoNormal><span style='mso-ansi-language:DE;mso-fareast-language:
		DE;mso-no-proof:yes'></span><span
		style='mso-ansi-language:DE'><o:p></o:p></span></p>
		</td>
		<td width=425 colspan=2 valign=top style='width:424.85pt;padding:7.2pt 5.75pt 2.9pt 5.75pt;
		height:37.8pt'>
		<h1><span style='mso-ansi-language:DE'>Rechnung<o:p></o:p></span></h1>
		</td>
		</tr>
		<tr style='mso-yfti-irow:1;height:23.4pt'>
		<w:Sdt ShowingPlcHdr="t" DocPart="13A398AD0E7B97469D3FA043010E4F38"
		ID="716560606">
		<td width=387 colspan=2 valign=top style='width:387.0pt;padding:0cm 5.75pt 2.9pt 5.75pt;
		height:23.4pt'>
		<p class=Slogan><span style='mso-ansi-language:DE'>[Ihr Firmenslogan]<o:p></o:p><w:sdtPr></w:sdtPr></span></p>
		</td>
		</w:Sdt>
		<td width=117 rowspan=2 valign=top style='width:117.0pt;padding:7.2pt 5.75pt 2.9pt 5.75pt;
		height:23.4pt'>
		<p class=DatumundNummer><span style='mso-ansi-language:DE'>Datum: <w:Sdt
		ShowingPlcHdr="t" DocPart="DF753200C3775348AE356AF9FAA0CE9B" Calendar="t"
		MapToDateTime="t" CalendarType="Gregorian" DateFormat="dd.MM.yyyy" Lang="DE"
		ID="716560632"><span class=MsoPlaceholderText>[Datum eingeben]</span></w:Sdt><o:p></o:p></span></p>
		<p class=DatumundNummer><span style='mso-ansi-language:DE'>RECHNUNG NR. <w:Sdt
		ShowingPlcHdr="t" DocPart="DEAC19F3DEF23940A41AA3E92731E343" Text="t"
		ID="716560634">[100]</w:Sdt><o:p></o:p></span></p>
		<p class=RechtsausgerichteterText><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		<w:Sdt ShowingPlcHdr="t" DocPart="AFACF1569441B345975887C1BB482942"
		ID="716560479">
		<p class=RechtsausgerichteterText><span class=MsoPlaceholderText><span
		style='mso-ansi-language:DE'>[Name]</span></span><span style='mso-ansi-language:
		DE'><o:p></o:p><w:sdtPr></w:sdtPr></span></p>
		</w:Sdt><w:Sdt ShowingPlcHdr="t" DocPart="8214BB5DEDB3634E899BEC1BE47C588F"
		ID="716560481">
		<p class=RechtsausgerichteterText><span style='mso-ansi-language:DE'>[Firmenname]<o:p></o:p><w:sdtPr></w:sdtPr></span></p>
		</w:Sdt><w:Sdt ShowingPlcHdr="t" DocPart="F79AE4723C6AB9419FE0299BC59315F8"
		ID="716560484">
		<p class=RechtsausgerichteterText><span style='mso-ansi-language:DE'>[Stra�e
		Hausnummer]<o:p></o:p><w:sdtPr></w:sdtPr></span></p>
		</w:Sdt><w:Sdt ShowingPlcHdr="t" DocPart="1F7AF04C930EA749BE0742DC565D21F5"
		ID="716560486">
		<p class=RechtsausgerichteterText><span class=MsoPlaceholderText><span
		style='mso-ansi-language:DE'>[Postleitzahl Ort]</span></span><span
		style='mso-ansi-language:DE'><o:p></o:p><w:sdtPr></w:sdtPr></span></p>
		</w:Sdt><w:Sdt ShowingPlcHdr="t" DocPart="36F6A57C8CE4D14A858A0FC9BB954817"
		ID="716560491">
		<p class=RechtsausgerichteterText><span style='mso-ansi-language:DE'>[Telefon]<o:p></o:p><w:sdtPr></w:sdtPr></span></p>
		</w:Sdt>
		<p class=RechtsausgerichteterText><span style='mso-ansi-language:DE'>Kundennummer
		<w:Sdt ShowingPlcHdr="t" DocPart="261BB3C27F56D548B82F61DC6664E4C4"
		ID="716560494">[ABC12345]</w:Sdt><o:p></o:p></span></p>
		</td>
		</tr>
		<tr style='mso-yfti-irow:2;mso-yfti-lastrow:yes;height:57.75pt'>
		<td width=387 colspan=2 valign=top style='width:387.0pt;padding:0cm 5.75pt 2.9pt 5.75pt;
		height:57.75pt'>
		<p class=berschriften><span style='mso-ansi-language:DE'>AN<o:p></o:p></span></p>
		</td>
		</tr>
		<![if !supportMisalignedColumns]>
		<tr height=0>
		<td width=80 style='border:none'></td>
		<td width=308 style='border:none'></td>
		<td width=117 style='border:none'></td>
		</tr>
		<![endif]>
		</table>
		
		</div>
		
		<p class=MsoNormal><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		
		<div align=center>
		
		<table class=MsoNormalTable border=1 cellspacing=0 cellpadding=0 width=504
		style='width:504.0pt;border-collapse:collapse;border:none;mso-border-alt:solid gray .5pt;
		mso-padding-alt:2.15pt 5.75pt 2.15pt 5.75pt;mso-border-insideh:.5pt solid gray;
		mso-border-insidev:.5pt solid gray'>
		<tr style='mso-yfti-irow:0;mso-yfti-firstrow:yes;page-break-inside:avoid;
		height:10.8pt'>
		<td width=62 style='width:62.45pt;border:solid #B0CCB0 1.0pt;mso-border-themecolor:
		accent2;mso-border-alt:solid #B0CCB0 .5pt;mso-border-themecolor:accent2;
		background:#EFF4EF;mso-background-themecolor:accent2;mso-background-themetint:
		51;padding:2.15pt 5.75pt 2.15pt 5.75pt;height:10.8pt'>
		<p class=Spaltenberschrift><span style='mso-ansi-language:DE'>Verk�ufer<o:p></o:p></span></p>
		</td>
		<td width=269 style='width:268.7pt;border:solid #B0CCB0 1.0pt;mso-border-themecolor:
		accent2;border-left:none;mso-border-left-alt:solid #B0CCB0 .5pt;mso-border-left-themecolor:
		accent2;mso-border-alt:solid #B0CCB0 .5pt;mso-border-themecolor:accent2;
		background:#EFF4EF;mso-background-themecolor:accent2;mso-background-themetint:
		51;padding:2.15pt 5.75pt 2.15pt 5.75pt;height:10.8pt'>
		<p class=Spaltenberschrift><span style='mso-ansi-language:DE'>Auftrag<o:p></o:p></span></p>
		</td>
		<td width=92 style='width:92.1pt;border:solid #B0CCB0 1.0pt;mso-border-themecolor:
		accent2;border-left:none;mso-border-left-alt:solid #B0CCB0 .5pt;mso-border-left-themecolor:
		accent2;mso-border-alt:solid #B0CCB0 .5pt;mso-border-themecolor:accent2;
		background:#EFF4EF;mso-background-themecolor:accent2;mso-background-themetint:
		51;padding:2.15pt 5.75pt 2.15pt 5.75pt;height:10.8pt'>
		<p class=Spaltenberschrift><span style='mso-ansi-language:DE'>Zahlungsbedingungen<o:p></o:p></span></p>
		</td>
		<td width=81 style='width:80.75pt;border:solid #B0CCB0 1.0pt;mso-border-themecolor:
		accent2;border-left:none;mso-border-left-alt:solid #B0CCB0 .5pt;mso-border-left-themecolor:
		accent2;mso-border-alt:solid #B0CCB0 .5pt;mso-border-themecolor:accent2;
		background:#EFF4EF;mso-background-themecolor:accent2;mso-background-themetint:
		51;padding:2.15pt 5.75pt 2.15pt 5.75pt;height:10.8pt'>
		<p class=Spaltenberschrift><span style='mso-ansi-language:DE'>F�llig am<o:p></o:p></span></p>
		</td>
		</tr>
		<tr style='mso-yfti-irow:1;mso-yfti-lastrow:yes;page-break-inside:avoid;
		height:10.8pt'>
		<td width=62 style='width:62.45pt;border:solid #B0CCB0 1.0pt;mso-border-themecolor:
		accent2;border-top:none;mso-border-top-alt:solid #B0CCB0 .5pt;mso-border-top-themecolor:
		accent2;mso-border-alt:solid #B0CCB0 .5pt;mso-border-themecolor:accent2;
		padding:2.15pt 5.75pt 2.15pt 5.75pt;height:10.8pt'>
		<p class=MsoNormal><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=269 style='width:268.7pt;border-top:none;border-left:none;
		border-bottom:solid #B0CCB0 1.0pt;mso-border-bottom-themecolor:accent2;
		border-right:solid #B0CCB0 1.0pt;mso-border-right-themecolor:accent2;
		mso-border-top-alt:solid #B0CCB0 .5pt;mso-border-top-themecolor:accent2;
		mso-border-left-alt:solid #B0CCB0 .5pt;mso-border-left-themecolor:accent2;
		mso-border-alt:solid #B0CCB0 .5pt;mso-border-themecolor:accent2;padding:2.15pt 5.75pt 2.15pt 5.75pt;
		height:10.8pt'>
		<p class=MsoNormal><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=92 style='width:92.1pt;border-top:none;border-left:none;border-bottom:
		solid #B0CCB0 1.0pt;mso-border-bottom-themecolor:accent2;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-top-alt:solid #B0CCB0 .5pt;
		mso-border-top-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-alt:solid #B0CCB0 .5pt;
		mso-border-themecolor:accent2;padding:2.15pt 5.75pt 2.15pt 5.75pt;height:
		10.8pt'>
		<p class=MsoNormal><span style='mso-ansi-language:DE'>F�llig bei Erhalt<o:p></o:p></span></p>
		</td>
		<td width=81 style='width:80.75pt;border-top:none;border-left:none;
		border-bottom:solid #B0CCB0 1.0pt;mso-border-bottom-themecolor:accent2;
		border-right:solid #B0CCB0 1.0pt;mso-border-right-themecolor:accent2;
		mso-border-top-alt:solid #B0CCB0 .5pt;mso-border-top-themecolor:accent2;
		mso-border-left-alt:solid #B0CCB0 .5pt;mso-border-left-themecolor:accent2;
		mso-border-alt:solid #B0CCB0 .5pt;mso-border-themecolor:accent2;padding:2.15pt 5.75pt 2.15pt 5.75pt;
		height:10.8pt'>
		<p class=MsoNormal><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		</tr>
		</table>
		
		</div>
		
		<p class=MsoNormal><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		
		<div align=center>
		
		<table class=MsoNormalTable border=1 cellspacing=0 cellpadding=0 width=504
		style='width:504.0pt;border-collapse:collapse;border:none;mso-border-alt:solid gray .5pt;
		mso-padding-alt:2.15pt 5.75pt 2.15pt 5.75pt;mso-border-insideh:.5pt solid gray;
		mso-border-insidev:.5pt solid gray'>
		<tr style='mso-yfti-irow:0;mso-yfti-firstrow:yes;page-break-inside:avoid;
		height:10.8pt'>
		<td width=63 style='width:63.0pt;border:solid #B0CCB0 1.0pt;mso-border-themecolor:
		accent2;mso-border-alt:solid #B0CCB0 .5pt;mso-border-themecolor:accent2;
		background:#EFF4EF;mso-background-themecolor:accent2;mso-background-themetint:
		51;padding:2.15pt 5.75pt 2.15pt 5.75pt;height:10.8pt'>
		<p class=Spaltenberschrift><span style='mso-ansi-language:DE'>Menge<o:p></o:p></span></p>
		</td>
		<td width=268 style='width:268.15pt;border:solid #B0CCB0 1.0pt;mso-border-themecolor:
		accent2;border-left:none;mso-border-left-alt:solid #B0CCB0 .5pt;mso-border-left-themecolor:
		accent2;mso-border-alt:solid #B0CCB0 .5pt;mso-border-themecolor:accent2;
		background:#EFF4EF;mso-background-themecolor:accent2;mso-background-themetint:
		51;padding:2.15pt 5.75pt 2.15pt 5.75pt;height:10.8pt'>
		<p class=Spaltenberschrift><span style='mso-ansi-language:DE'>Beschreibung</span><span
		style='mso-ansi-language:DE;mso-fareast-language:ZH-CN'><o:p></o:p></span></p>
		</td>
		<td width=92 style='width:91.85pt;border:solid #B0CCB0 1.0pt;mso-border-themecolor:
		accent2;border-left:none;mso-border-left-alt:solid #B0CCB0 .5pt;mso-border-left-themecolor:
		accent2;mso-border-alt:solid #B0CCB0 .5pt;mso-border-themecolor:accent2;
		background:#EFF4EF;mso-background-themecolor:accent2;mso-background-themetint:
		51;padding:2.15pt 5.75pt 2.15pt 5.75pt;height:10.8pt'>
		<p class=Spaltenberschrift><span style='mso-ansi-language:DE'>Preis pro
		Einheit<o:p></o:p></span></p>
		</td>
		<td width=81 style='width:81.0pt;border:solid #B0CCB0 1.0pt;mso-border-themecolor:
		accent2;border-left:none;mso-border-left-alt:solid #B0CCB0 .5pt;mso-border-left-themecolor:
		accent2;mso-border-alt:solid #B0CCB0 .5pt;mso-border-themecolor:accent2;
		background:#EFF4EF;mso-background-themecolor:accent2;mso-background-themetint:
		51;padding:2.15pt 5.75pt 2.15pt 5.75pt;height:10.8pt'>
		<p class=Spaltenberschrift><span style='mso-ansi-language:DE'>Summe der
		Positionen<o:p></o:p></span></p>
		</td>
		</tr>
		<tr style='mso-yfti-irow:1;page-break-inside:avoid;height:10.8pt'>
		<td width=63 style='width:63.0pt;border-top:none;border-left:solid #B0CCB0 1.0pt;
		mso-border-left-themecolor:accent2;border-bottom:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-top-alt:solid #B0CCB0 .5pt;
		mso-border-top-themecolor:accent2;mso-border-top-alt:solid #B0CCB0 .5pt;
		mso-border-top-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 5.75pt 2.15pt 5.75pt;
		height:10.8pt'>
		<p class=MsoNormal><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=268 style='width:268.15pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-top-alt:solid #B0CCB0 .5pt;
		mso-border-top-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-top-alt:solid #B0CCB0 .5pt;
		mso-border-top-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 5.75pt 2.15pt 5.75pt;
		height:10.8pt'>
		<p class=MsoNormal><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=92 style='width:91.85pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-top-alt:solid #B0CCB0 .5pt;
		mso-border-top-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-top-alt:solid #B0CCB0 .5pt;
		mso-border-top-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 10.8pt 2.15pt 10.8pt;
		height:10.8pt'>
		<p class=Betrag><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=81 style='width:81.0pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-top-alt:solid #B0CCB0 .5pt;
		mso-border-top-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-top-alt:solid #B0CCB0 .5pt;
		mso-border-top-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 10.8pt 2.15pt 10.8pt;
		height:10.8pt'>
		<p class=Betrag><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		</tr>
		<tr style='mso-yfti-irow:2;page-break-inside:avoid;height:10.8pt'>
		<td width=63 style='width:63.0pt;border-top:none;border-left:solid #B0CCB0 1.0pt;
		mso-border-left-themecolor:accent2;border-bottom:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 5.75pt 2.15pt 5.75pt;
		height:10.8pt'>
		<p class=MsoNormal><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=268 style='width:268.15pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 5.75pt 2.15pt 5.75pt;
		height:10.8pt'>
		<p class=MsoNormal><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=92 style='width:91.85pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 10.8pt 2.15pt 10.8pt;
		height:10.8pt'>
		<p class=Betrag><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=81 style='width:81.0pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 10.8pt 2.15pt 10.8pt;
		height:10.8pt'>
		<p class=Betrag><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		</tr>
		<tr style='mso-yfti-irow:3;page-break-inside:avoid;height:10.8pt'>
		<td width=63 style='width:63.0pt;border-top:none;border-left:solid #B0CCB0 1.0pt;
		mso-border-left-themecolor:accent2;border-bottom:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 5.75pt 2.15pt 5.75pt;
		height:10.8pt'>
		<p class=MsoNormal><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=268 style='width:268.15pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 5.75pt 2.15pt 5.75pt;
		height:10.8pt'>
		<p class=MsoNormal><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=92 style='width:91.85pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 10.8pt 2.15pt 10.8pt;
		height:10.8pt'>
		<p class=Betrag><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=81 style='width:81.0pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 10.8pt 2.15pt 10.8pt;
		height:10.8pt'>
		<p class=Betrag><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		</tr>
		<tr style='mso-yfti-irow:4;page-break-inside:avoid;height:10.8pt'>
		<td width=63 style='width:63.0pt;border-top:none;border-left:solid #B0CCB0 1.0pt;
		mso-border-left-themecolor:accent2;border-bottom:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 5.75pt 2.15pt 5.75pt;
		height:10.8pt'>
		<p class=MsoNormal><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=268 style='width:268.15pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 5.75pt 2.15pt 5.75pt;
		height:10.8pt'>
		<p class=MsoNormal><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=92 style='width:91.85pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 10.8pt 2.15pt 10.8pt;
		height:10.8pt'>
		<p class=Betrag><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=81 style='width:81.0pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 10.8pt 2.15pt 10.8pt;
		height:10.8pt'>
		<p class=Betrag><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		</tr>
		<tr style='mso-yfti-irow:5;page-break-inside:avoid;height:10.8pt'>
		<td width=63 style='width:63.0pt;border-top:none;border-left:solid #B0CCB0 1.0pt;
		mso-border-left-themecolor:accent2;border-bottom:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 5.75pt 2.15pt 5.75pt;
		height:10.8pt'>
		<p class=MsoNormal><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=268 style='width:268.15pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 5.75pt 2.15pt 5.75pt;
		height:10.8pt'>
		<p class=MsoNormal><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=92 style='width:91.85pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 10.8pt 2.15pt 10.8pt;
		height:10.8pt'>
		<p class=Betrag><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=81 style='width:81.0pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 10.8pt 2.15pt 10.8pt;
		height:10.8pt'>
		<p class=Betrag><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		</tr>
		<tr style='mso-yfti-irow:6;page-break-inside:avoid;height:10.8pt'>
		<td width=63 style='width:63.0pt;border-top:none;border-left:solid #B0CCB0 1.0pt;
		mso-border-left-themecolor:accent2;border-bottom:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 5.75pt 2.15pt 5.75pt;
		height:10.8pt'>
		<p class=MsoNormal><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=268 style='width:268.15pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 5.75pt 2.15pt 5.75pt;
		height:10.8pt'>
		<p class=MsoNormal><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=92 style='width:91.85pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 10.8pt 2.15pt 10.8pt;
		height:10.8pt'>
		<p class=Betrag><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=81 style='width:81.0pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 10.8pt 2.15pt 10.8pt;
		height:10.8pt'>
		<p class=Betrag><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		</tr>
		<tr style='mso-yfti-irow:7;page-break-inside:avoid;height:10.8pt'>
		<td width=63 style='width:63.0pt;border-top:none;border-left:solid #B0CCB0 1.0pt;
		mso-border-left-themecolor:accent2;border-bottom:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 5.75pt 2.15pt 5.75pt;
		height:10.8pt'>
		<p class=MsoNormal><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=268 style='width:268.15pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 5.75pt 2.15pt 5.75pt;
		height:10.8pt'>
		<p class=MsoNormal><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=92 style='width:91.85pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 10.8pt 2.15pt 10.8pt;
		height:10.8pt'>
		<p class=Betrag><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=81 style='width:81.0pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 10.8pt 2.15pt 10.8pt;
		height:10.8pt'>
		<p class=Betrag><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		</tr>
		<tr style='mso-yfti-irow:8;page-break-inside:avoid;height:10.8pt'>
		<td width=63 style='width:63.0pt;border-top:none;border-left:solid #B0CCB0 1.0pt;
		mso-border-left-themecolor:accent2;border-bottom:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 5.75pt 2.15pt 5.75pt;
		height:10.8pt'>
		<p class=MsoNormal><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=268 style='width:268.15pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 5.75pt 2.15pt 5.75pt;
		height:10.8pt'>
		<p class=MsoNormal><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=92 style='width:91.85pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 10.8pt 2.15pt 10.8pt;
		height:10.8pt'>
		<p class=Betrag><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=81 style='width:81.0pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 10.8pt 2.15pt 10.8pt;
		height:10.8pt'>
		<p class=Betrag><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		</tr>
		<tr style='mso-yfti-irow:9;page-break-inside:avoid;height:10.8pt'>
		<td width=63 style='width:63.0pt;border-top:none;border-left:solid #B0CCB0 1.0pt;
		mso-border-left-themecolor:accent2;border-bottom:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 5.75pt 2.15pt 5.75pt;
		height:10.8pt'>
		<p class=MsoNormal><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=268 style='width:268.15pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 5.75pt 2.15pt 5.75pt;
		height:10.8pt'>
		<p class=MsoNormal><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=92 style='width:91.85pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 10.8pt 2.15pt 10.8pt;
		height:10.8pt'>
		<p class=Betrag><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=81 style='width:81.0pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 10.8pt 2.15pt 10.8pt;
		height:10.8pt'>
		<p class=Betrag><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		</tr>
		<tr style='mso-yfti-irow:10;page-break-inside:avoid;height:10.8pt'>
		<td width=63 style='width:63.0pt;border-top:none;border-left:solid #B0CCB0 1.0pt;
		mso-border-left-themecolor:accent2;border-bottom:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 5.75pt 2.15pt 5.75pt;
		height:10.8pt'>
		<p class=MsoNormal><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=268 style='width:268.15pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 5.75pt 2.15pt 5.75pt;
		height:10.8pt'>
		<p class=MsoNormal><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=92 style='width:91.85pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 10.8pt 2.15pt 10.8pt;
		height:10.8pt'>
		<p class=Betrag><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=81 style='width:81.0pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 10.8pt 2.15pt 10.8pt;
		height:10.8pt'>
		<p class=Betrag><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		</tr>
		<tr style='mso-yfti-irow:11;page-break-inside:avoid;height:10.8pt'>
		<td width=63 style='width:63.0pt;border-top:none;border-left:solid #B0CCB0 1.0pt;
		mso-border-left-themecolor:accent2;border-bottom:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 5.75pt 2.15pt 5.75pt;
		height:10.8pt'>
		<p class=MsoNormal><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=268 style='width:268.15pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 5.75pt 2.15pt 5.75pt;
		height:10.8pt'>
		<p class=MsoNormal><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=92 style='width:91.85pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 10.8pt 2.15pt 10.8pt;
		height:10.8pt'>
		<p class=Betrag><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=81 style='width:81.0pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 10.8pt 2.15pt 10.8pt;
		height:10.8pt'>
		<p class=Betrag><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		</tr>
		<tr style='mso-yfti-irow:12;page-break-inside:avoid;height:10.8pt'>
		<td width=63 style='width:63.0pt;border-top:none;border-left:solid #B0CCB0 1.0pt;
		mso-border-left-themecolor:accent2;border-bottom:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 5.75pt 2.15pt 5.75pt;
		height:10.8pt'>
		<p class=MsoNormal><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=268 style='width:268.15pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 5.75pt 2.15pt 5.75pt;
		height:10.8pt'>
		<p class=MsoNormal><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=92 style='width:91.85pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 10.8pt 2.15pt 10.8pt;
		height:10.8pt'>
		<p class=Betrag><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=81 style='width:81.0pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 10.8pt 2.15pt 10.8pt;
		height:10.8pt'>
		<p class=Betrag><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		</tr>
		<tr style='mso-yfti-irow:13;page-break-inside:avoid;height:10.8pt'>
		<td width=63 style='width:63.0pt;border-top:none;border-left:solid #B0CCB0 1.0pt;
		mso-border-left-themecolor:accent2;border-bottom:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 5.75pt 2.15pt 5.75pt;
		height:10.8pt'>
		<p class=MsoNormal><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=268 style='width:268.15pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 5.75pt 2.15pt 5.75pt;
		height:10.8pt'>
		<p class=MsoNormal><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=92 style='width:91.85pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 10.8pt 2.15pt 10.8pt;
		height:10.8pt'>
		<p class=Betrag><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=81 style='width:81.0pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 10.8pt 2.15pt 10.8pt;
		height:10.8pt'>
		<p class=Betrag><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		</tr>
		<tr style='mso-yfti-irow:14;page-break-inside:avoid;height:10.8pt'>
		<td width=63 style='width:63.0pt;border-top:none;border-left:solid #B0CCB0 1.0pt;
		mso-border-left-themecolor:accent2;border-bottom:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 5.75pt 2.15pt 5.75pt;
		height:10.8pt'>
		<p class=MsoNormal><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=268 style='width:268.15pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 5.75pt 2.15pt 5.75pt;
		height:10.8pt'>
		<p class=MsoNormal><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=92 style='width:91.85pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 10.8pt 2.15pt 10.8pt;
		height:10.8pt'>
		<p class=Betrag><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=81 style='width:81.0pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 10.8pt 2.15pt 10.8pt;
		height:10.8pt'>
		<p class=Betrag><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		</tr>
		<tr style='mso-yfti-irow:15;page-break-inside:avoid;height:10.8pt'>
		<td width=63 style='width:63.0pt;border-top:none;border-left:solid #B0CCB0 1.0pt;
		mso-border-left-themecolor:accent2;border-bottom:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 5.75pt 2.15pt 5.75pt;
		height:10.8pt'>
		<p class=MsoNormal><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=268 style='width:268.15pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 5.75pt 2.15pt 5.75pt;
		height:10.8pt'>
		<p class=MsoNormal><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=92 style='width:91.85pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 10.8pt 2.15pt 10.8pt;
		height:10.8pt'>
		<p class=Betrag><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=81 style='width:81.0pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 10.8pt 2.15pt 10.8pt;
		height:10.8pt'>
		<p class=Betrag><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		</tr>
		<tr style='mso-yfti-irow:16;page-break-inside:avoid;height:10.8pt'>
		<td width=63 style='width:63.0pt;border:solid #B0CCB0 1.0pt;mso-border-themecolor:
		accent2;border-top:none;mso-border-left-alt:solid #B0CCB0 .5pt;mso-border-left-themecolor:
		accent2;mso-border-bottom-alt:solid #B0CCB0 .5pt;mso-border-bottom-themecolor:
		accent2;mso-border-right-alt:solid #B0CCB0 .5pt;mso-border-right-themecolor:
		accent2;padding:2.15pt 5.75pt 2.15pt 5.75pt;height:10.8pt'>
		<p class=MsoNormal><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=268 style='width:268.15pt;border-top:none;border-left:none;
		border-bottom:solid #B0CCB0 1.0pt;mso-border-bottom-themecolor:accent2;
		border-right:solid #B0CCB0 1.0pt;mso-border-right-themecolor:accent2;
		mso-border-left-alt:solid #B0CCB0 .5pt;mso-border-left-themecolor:accent2;
		mso-border-left-alt:solid #B0CCB0 .5pt;mso-border-left-themecolor:accent2;
		mso-border-bottom-alt:solid #B0CCB0 .5pt;mso-border-bottom-themecolor:accent2;
		mso-border-right-alt:solid #B0CCB0 .5pt;mso-border-right-themecolor:accent2;
		padding:2.15pt 5.75pt 2.15pt 5.75pt;height:10.8pt'>
		<p class=MsoNormal><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=92 style='width:91.85pt;border-top:none;border-left:none;
		border-bottom:solid #B0CCB0 1.0pt;mso-border-bottom-themecolor:accent2;
		border-right:solid #B0CCB0 1.0pt;mso-border-right-themecolor:accent2;
		mso-border-left-alt:solid #B0CCB0 .5pt;mso-border-left-themecolor:accent2;
		mso-border-left-alt:solid #B0CCB0 .5pt;mso-border-left-themecolor:accent2;
		mso-border-bottom-alt:solid #B0CCB0 .5pt;mso-border-bottom-themecolor:accent2;
		mso-border-right-alt:solid #B0CCB0 .5pt;mso-border-right-themecolor:accent2;
		padding:2.15pt 10.8pt 2.15pt 10.8pt;height:10.8pt'>
		<p class=Betrag><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=81 style='width:81.0pt;border-top:none;border-left:none;border-bottom:
		solid #B0CCB0 1.0pt;mso-border-bottom-themecolor:accent2;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-bottom-alt:solid #B0CCB0 .5pt;
		mso-border-bottom-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 10.8pt 2.15pt 10.8pt;
		height:10.8pt'>
		<p class=Betrag><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		</tr>
		<tr style='mso-yfti-irow:17;page-break-inside:avoid;height:10.8pt'>
		<td width=331 colspan=2 rowspan=3 style='width:331.15pt;border:none;
		mso-border-top-alt:solid #B0CCB0 .5pt;mso-border-top-themecolor:accent2;
		padding:2.15pt 5.75pt 2.15pt 5.75pt;height:10.8pt'>
		<p class=MsoNormal><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		<td width=92 style='width:91.85pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-top-alt:solid #B0CCB0 .5pt;
		mso-border-top-themecolor:accent2;mso-border-top-alt:solid #B0CCB0 .5pt;
		mso-border-top-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 5.75pt 2.15pt 5.75pt;
		height:10.8pt'>
		<p class=RechtsausgerichteterText><span style='mso-ansi-language:DE'>Zwischensumme<o:p></o:p></span></p>
		</td>
		<td width=81 style='width:81.0pt;border-top:none;border-left:none;border-bottom:
		solid #B0CCB0 1.0pt;mso-border-bottom-themecolor:accent2;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-top-alt:solid #B0CCB0 .5pt;
		mso-border-top-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-alt:solid #B0CCB0 .5pt;
		mso-border-themecolor:accent2;padding:2.15pt 10.8pt 2.15pt 10.8pt;height:
		10.8pt'>
		<p class=Betrag><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		</tr>
		<tr style='mso-yfti-irow:18;page-break-inside:avoid;height:10.8pt'>
		<td width=92 style='width:91.85pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 5.75pt 2.15pt 5.75pt;
		height:10.8pt'>
		<p class=RechtsausgerichteterText><span style='mso-ansi-language:DE'>MwSt<o:p></o:p></span></p>
		</td>
		<td width=81 style='width:81.0pt;border-top:none;border-left:none;border-bottom:
		solid #B0CCB0 1.0pt;mso-border-bottom-themecolor:accent2;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-top-alt:solid #B0CCB0 .5pt;
		mso-border-top-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-alt:solid #B0CCB0 .5pt;
		mso-border-themecolor:accent2;padding:2.15pt 10.8pt 2.15pt 10.8pt;height:
		10.8pt'>
		<p class=Betrag><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		</tr>
		<tr style='mso-yfti-irow:19;mso-yfti-lastrow:yes;page-break-inside:avoid;
		height:10.8pt'>
		<td width=92 style='width:91.85pt;border:none;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-right-alt:solid #B0CCB0 .5pt;
		mso-border-right-themecolor:accent2;padding:2.15pt 5.75pt 2.15pt 5.75pt;
		height:10.8pt'>
		<p class=RechtsausgerichteterText><span style='mso-ansi-language:DE'>Summe<o:p></o:p></span></p>
		</td>
		<td width=81 style='width:81.0pt;border-top:none;border-left:none;border-bottom:
		solid #B0CCB0 1.0pt;mso-border-bottom-themecolor:accent2;border-right:solid #B0CCB0 1.0pt;
		mso-border-right-themecolor:accent2;mso-border-top-alt:solid #B0CCB0 .5pt;
		mso-border-top-themecolor:accent2;mso-border-left-alt:solid #B0CCB0 .5pt;
		mso-border-left-themecolor:accent2;mso-border-alt:solid #B0CCB0 .5pt;
		mso-border-themecolor:accent2;padding:2.15pt 10.8pt 2.15pt 10.8pt;height:
		10.8pt'>
		<p class=Betrag><span style='mso-ansi-language:DE'><o:p>&nbsp;</o:p></span></p>
		</td>
		</tr>
		</table>
		
		</div>
		
		<p class=ZentrierterTextunten><span style='mso-ansi-language:DE'>Stellen Sie
		alle �berweisungen auf <w:Sdt ShowingPlcHdr="t"
		DocPart="DB677D5F5BC0174B8CD87730C9C748BC" ID="716560639">[Ihr Firmenname]</w:Sdt></span><span
		style='mso-ansi-language:DE;mso-fareast-language:ZH-CN'> </span><span
		style='mso-ansi-language:DE'>als Empf�nger aus.<o:p></o:p></span></p>
		
		<p class=VielenDank><span style='mso-ansi-language:DE'>Vielen Dank f�r Ihre
		Bestellung!<o:p></o:p></span></p>
		
		<p class=ZentrierterTextunten><span style='mso-ansi-language:DE'><w:Sdt
		ShowingPlcHdr="t" DocPart="D60455278C981E4C878845B7F63973B5" ID="716560525">[Ihr
		Firmenname]</w:Sdt> <w:Sdt ShowingPlcHdr="t"
		DocPart="FE721DC6DC63C94BBB9116828E0C679D" ID="716560527">[Stra�e Hausnummer]</w:Sdt>
		<w:Sdt ShowingPlcHdr="t" DocPart="05BEB0E5ED77744B8DBF130575325817"
		ID="716560530">[Postleitzahl Ort]</w:Sdt> Telefon <w:Sdt ShowingPlcHdr="t"
		DocPart="65A838C819567E4DB09F7B5152331E91" ID="716560532">[000-000-0000]</w:Sdt>
		Fax <w:Sdt ShowingPlcHdr="t" DocPart="3EF0F9D75C90F34F983E607B69D16ACA"
		ID="716560539">[000-000-0000]</w:Sdt> <w:Sdt ShowingPlcHdr="t"
		DocPart="BD3C10482759F743A297442EE00FEE55" ID="716560542">[E-Mail-Adresse]</w:Sdt><o:p></o:p></span></p>
		
		</div>	
		</body>
	`

	html := pdf.HTMLBasicNew()
	html.Write(lineHt, htmlStr)
	return pdf.Output(writer)
}
