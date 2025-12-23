package domain

type Payload struct {
	MerchantName       string        `json:"merchant_name"`
	NoOrder            string        `json:"no_order"`
	Tanggal            string        `json:"tanggal"`
	TypeTransaction    string        `json:"type_transaction"`
	Servis             string        `json:"servis"`
	KasirName          string        `json:"kasir_name"`
	CustomerID         string        `json:"customer_id"`
	CustomerName       string        `json:"customer_name"`
	CustomerPhone      string        `json:"customer_phone"`
	CustomerAddress    string        `json:"customer_address"`
	Notes              string        `json:"notes"`
	MetodePembayaran   string        `json:"metode_pembayaran"`
	SubTotal           int64         `json:"sub_total"`
	Ongkir             int64         `json:"ongkir"`
	Pajak              int64         `json:"pajak"`
	Diskon             int64         `json:"diskon"`
	Total              int64         `json:"total"`
	MoneyReceived      int64         `json:"money_received"`
	MoneyBack          int64         `json:"money_back"`
	DetailOrder        []DetailOrder `json:"detail_order"`
	Logo               string        `json:"logo"`
	AdvanceHeader      bool          `json:"advance_header"`
	SummaryPayment     []any         `json:"summary_payment"`
	Header             Header        `json:"header"`
	Footer             Footer        `json:"footer"`
	TransactionPurpose string        `json:"transaction_purpose"`
	SourceOfFund       string        `json:"source_of_fund"`
}

type Header struct {
	Title                 string `json:"title"`
	Address               string `json:"address"`
	Phone                 string `json:"phone"`
	BusinessLicenseNumber string `json:"business_license_number"`
	BusinessLicenseDate   string `json:"business_license_date"`
}

type Footer struct {
	Statement string `json:"statement"`
	NoteID    string `json:"note_id"`
	NoteEN    string `json:"note_en"`
}

type DetailOrder struct {
	Currency string      `json:"currency"`
	Amount   interface{} `json:"amount"`
	Rate     string      `json:"rate"`
	Total    interface{} `json:"total"`
	Denom    string      `json:"denom"`
}

type LineWriter interface {
	AppendRaw(string)
	NewLine()
}

type Settings struct {
	IDMerchant   string
	NameMerchant string
	PrinterName  string
	PaperWidth   int
	Offsite      int
	TopicPrefix  string
}

type SettingsRepository interface {
	Get() (Settings, error)
}
