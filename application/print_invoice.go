package application

import (
	"other/go-printer-termal/domain"
	"strings"
)

type PrintInvoiceUseCase struct {
	PaperWidth int
}

func (uc PrintInvoiceUseCase) Execute(p PrinterEngine, payload domain.Payload) error {
	if payload.Logo != "" {
		_ = p.PrintImageURLScaled(payload.Logo, 200)
	}

	p.AlignCenter()
	p.Size(2, 2)
	p.BoldOn()
	p.AppendRaw(payload.MerchantName + "\n")
	p.SizeReset()
	p.BoldOff()

	p.AlignCenter()
	p.AppendRaw("Authorized Money Changer \n")
	p.SizeReset()
	p.BoldOff()

	p.AlignCenter()
	p.AppendRaw("No Izin: " + payload.Header.BusinessLicenseNumber + "\n")
	p.SizeReset()
	p.BoldOff()

	p.AlignCenter()
	p.AppendRaw("Tgl Terbit: " + payload.Header.BusinessLicenseDate + "\n")
	p.SizeReset()
	p.BoldOff()

	if payload.AdvanceHeader {
		p.AlignCenter()
		domain.PrintLeftAlignedText(p, uc.PaperWidth, payload.Header.Address)
		p.AppendRaw(payload.Header.Phone + "\n")
		p.NewLines(1)
	} else {
		p.NewLines(1)
	}

	p.AlignCenter()
	p.Size(2, 2)
	p.BoldOn()
	p.AppendRaw(payload.Header.Title + "\n")
	p.SizeReset()
	p.BoldOff()
	p.NewLine()

	p.AlignLeft()
	domain.PrintInvoiceItem(p, uc.PaperWidth, "No. Transaksi    :", payload.NoOrder)
	domain.PrintInvoiceItem(p, uc.PaperWidth, "Tgl. Transaksi   :", payload.Tanggal)
	domain.PrintInvoiceItem(p, uc.PaperWidth, "Pelanggan        :", payload.CustomerName)
	domain.PrintInvoiceItem(p, uc.PaperWidth, "No CIF           :", payload.CustomerID)
	domain.PrintInvoiceItem(p, uc.PaperWidth, "No ID            :", "")
	p.AppendRaw("Alamat\n")
	domain.PrintLeftAlignedText(p, uc.PaperWidth, payload.CustomerAddress)
	domain.PrintInvoiceItem(p, uc.PaperWidth, "No. Telp         :", payload.CustomerPhone)
	domain.PrintInvoiceItem(p, uc.PaperWidth, "Tujuan Transaksi :", "")
	domain.PrintInvoiceItem(p, uc.PaperWidth, "Catatan          :", payload.Notes)
	p.NewLines(2)

	p.AlignLeft()
	p.Size(1, 1)
	p.BoldOn()
	p.AppendRaw("Detail Pesanan\n")
	p.BoldOff()
	p.AppendRaw(domain.RepeatRune('-', uc.PaperWidth) + "\n")

	for _, d := range payload.DetailOrder {
		amountStr := domain.FormatNumberID(int64(d.Amount.(float64)))
		totalStr := domain.FormatNumberID(int64(d.Total.(float64)))
		domain.PrintInvoiceItem(p, uc.PaperWidth, "Currency :", d.Currency)
		domain.PrintInvoiceItem(p, uc.PaperWidth, "Amount   :", amountStr)
		domain.PrintInvoiceItem(p, uc.PaperWidth, "Rate     :", d.Rate)
		domain.PrintInvoiceItem(p, uc.PaperWidth, "Total    :", "Rp "+totalStr)
		p.AppendRaw(domain.RepeatRune('-', uc.PaperWidth) + "\n")
	}
	p.NewLine()

	domain.PrintInvoiceItem(p, uc.PaperWidth, "Sub Total :", "Rp "+domain.FormatNumberID(payload.SubTotal))
	if payload.Ongkir != 0 {
		domain.PrintInvoiceItem(p, uc.PaperWidth, "Ongkir    :", "Rp "+domain.FormatNumberID(payload.Ongkir))
	}
	if payload.Pajak != 0 {
		domain.PrintInvoiceItem(p, uc.PaperWidth, "Pajak     :", "Rp "+domain.FormatNumberID(payload.Pajak))
	}
	if payload.Diskon != 0 {
		domain.PrintInvoiceItem(p, uc.PaperWidth, "Diskon    :", "Rp "+domain.FormatNumberID(payload.Diskon))
	}
	domain.PrintInvoiceItem(p, uc.PaperWidth, "Total     :", "Rp "+domain.FormatNumberID(payload.Total))
	p.AppendRaw(domain.RepeatRune('-', uc.PaperWidth) + "\n")

	p.BoldOn()
	p.AppendRaw("Pembayaran\n")
	p.BoldOff()
	p.AlignLeft()
	p.AppendRaw("Metode Pembayaran :\n")
	domain.PrintLeftAlignedText(p, uc.PaperWidth, payload.MetodePembayaran)
	p.NewLine()
	p.AppendRaw(domain.RepeatRune('-', uc.PaperWidth) + "\n")

	domain.PrintInvoiceItem(p, uc.PaperWidth, "Pembayaran :", "Rp "+domain.FormatNumberID(payload.MoneyReceived))
	domain.PrintInvoiceItem(p, uc.PaperWidth, "Kembalian  :", "Rp "+domain.FormatNumberID(payload.MoneyBack))
	p.NewLines(2)

	if payload.Footer.Statement != "" {
		domain.PrintLeftAlignedText(p, uc.PaperWidth, payload.Footer.Statement)
		p.NewLines(1)
	}
	domain.PrintLeftAlignedText(p, uc.PaperWidth, strings.ReplaceAll(payload.Footer.NoteID, "<br>", "\n"))
	p.NewLines(2)

	kasir := domain.ProcessName(payload.KasirName)
	customer := domain.ProcessName(payload.CustomerName)
	domain.PrintInvoiceItem(p, uc.PaperWidth, "    Served By :", "Customer    ")
	p.NewLines(6)
	domain.PrintInvoiceItem(p, uc.PaperWidth, kasir, customer)
	p.NewLines(2)

	p.BoldOn()
	p.AppendRaw("Perhatian / Note\n")
	p.BoldOff()
	domain.PrintLeftAlignedText(p, uc.PaperWidth, strings.ReplaceAll(payload.Footer.NoteID, "<br>", "\n"))
	domain.PrintLeftAlignedText(p, uc.PaperWidth, strings.ReplaceAll(payload.Footer.NoteEN, "<br>", "\n"))
	p.NewLines(6)
	p.SizeReset()
	return nil
}
