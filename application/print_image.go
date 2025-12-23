package application

import "other/go-printer-termal/domain"

func (uc PrintInvoiceUseCase) ExecuteTestImageTest(p PrinterEngine, payload domain.Payload) error {
	_ = p.PrintImageURLScaled(payload.Logo, 200)
	return nil
}
