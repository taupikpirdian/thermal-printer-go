package main

import (
	"encoding/json"
	"log"
	"other/go-printer-termal/application"
	"other/go-printer-termal/domain"
	escpos "other/go-printer-termal/infrastructure/escpos"
	"strings"

	"github.com/alexbrainman/printer"
)

func main() {
	printerName := "EPSON TM-U220 Receipt"
	paperWidth := 42

	jsonPayload := `{"merchant_name":"Midas xchange","no_order":"FB202511147464","tanggal":"2025-11-14","type_transaction":"purchase","servis":"Pembelian","kasir_name":"Kasir Midas","customer_id":"CS202511146494","customer_name":"SHEN HSUEH CHENG","customer_phone":"-","customer_address":"TAIWAN","notes":"-","metode_pembayaran":"Cash","sub_total":1805000,"ongkir":0,"pajak":0,"diskon":0,"total":1805000,"money_received":1805000,"money_back":0,"detail_order":[{"currency":"USD","amount":100,"rate":"16.550","total":1655000,"denom":"100"},{"currency":"USD","amount":10,"rate":"15.000","total":150000,"denom":"1"}],"logo":"https://dev-dashboard.valast.id/file/merchant/36dbaa22-4628-4ecd-968a-95b739535e0d/Ohcw6mrkZQ-20241209142948.jpg","advance_header":false,"summary_payment":[],"header":{"title":"INVOICE PEMBELIAN","address":"AEON Mall BSD Jl. BSD Raya Utama Lt Dasar","phone":"081316348682","business_license_number":"TEST/001","business_license_date":"13 November 2025"},"footer":{"statement":"","note_id":"Kekurangan penerimaan uang tidak ditanggung setelah meninggalkan loket/kasir.<br>Dokument pendukung yang diberikan kepada Kami untuk kepentingan pembelian valas adalah asli dan benar adanya.","note_en":"Deficiancies in cash receipts are not covered after leaving the counter/cashier.<br>Supporting documents provided to us for the purpose of purchasing foreign exchange are original and true."}}`

	var payload domain.Payload
	if err := json.Unmarshal([]byte(jsonPayload), &payload); err != nil {
		log.Fatalf("parse payload: %v", err)
	}

	p, err := printer.Open(printerName)
	if err != nil {
		log.Fatalf("open printer %s: %v", printerName, err)
	}
	defer p.Close()

	if err = p.StartRawDocument("Receipt"); err != nil {
		log.Fatalf("StartRawDocument: %v", err)
	}
	if err = p.StartPage(); err != nil {
		log.Fatalf("StartPage: %v", err)
	}

	ep := escpos.NewEscposPrinter(p)
	uc := application.PrintInvoiceUseCase{PaperWidth: paperWidth}

	if err = uc.ExecuteTestImage(ep, payload); err != nil {
		log.Fatalf("print: %v", err)
	}

	ep.Cut()
	if err = ep.PrintAndFlush(); err != nil {
		log.Fatalf("flush: %v", err)
	}

	if err = p.EndPage(); err != nil {
		log.Printf("EndPage: %v", err)
	}
	if err = p.EndDocument(); err != nil {
		log.Printf("EndDocument: %v", err)
	}

	log.Println(strings.TrimSpace("Pencetakan Selesai"))
}
