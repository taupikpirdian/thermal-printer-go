package application

import "other/go-printer-termal/domain"

type PrinterEngine interface {
	AlignLeft()
	AlignCenter()
	AlignRight()
	BoldOn()
	BoldOff()
	UnderlineOn()
	UnderlineOff()
	Size(width, height int)
	SizeReset()
	NewLine()
	NewLines(n int)
	AppendRaw(s string)
	Cut()
	PrintAndFlush() error
	PrintImageURL(url string) error
	PrintImageURLScaled(url string, maxWidth int) error
}

type InvoiceUseCase interface {
	Execute(p PrinterEngine, payload domain.Payload) error
	ExecuteTestImageTest(p PrinterEngine, payload domain.Payload) error
}
