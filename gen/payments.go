package gen

// PaymentRow is a row from the PAYMENTS table
type PaymentRow struct {
	RegionID           string `db:"-" csv:"REGION_ID"`
	PaymentType        string `db:"-" csv:"PAYMENT_TYPE"`
	PayTypeID          string `db:"method" csv:"PAY_TYPE_ID"`
	PaymentDate        string `db:"date" csv:"PAYMENT_DATE"`
	Amount             string `db:"sum" csv:"AMOUNT"`
	AmountCurrency     string `db:"-" csv:"AMOUNT_CURRENCY"`
	PhoneNumber        string `db:"-" csv:"PHONE_NUMBER"`
	Account            string `db:"-" csv:"ACCOUNT"`
	AbonentID          string `db:"uid" csv:"ABONENT_ID"`
	BankAccount        string `db:"-" csv:"BANK_ACCOUNT"`
	BankName           string `db:"-" csv:"BANK_NAME"`
	ExpressCardNumber  string `db:"-" csv:"EXPRESS_CARD_NUMBER"`
	TerminalID         string `db:"-" csv:"TERMINAL_ID"`
	TerminalNumber     string `db:"-" csv:"TERMINAL_NUMBER"`
	LATITUDE           string `db:"-" csv:"LATITUDE"`
	LONGITUDE          string `db:"-" csv:"LONGITUDE"`
	ProjectionType     string `db:"-" csv:"PROJECTION_TYPE"`
	CenterID           string `db:"-" csv:"CENTER_ID"`
	DonatedPhoneNumber string `db:"-" csv:"DONATED_PHONE_NUMBER"`
	DonatedAccount     string `db:"-" csv:"DONATED_ACCOUNT"`
	DonatedInternalID1 string `db:"-" csv:"DONATED_INTERNAL_ID1"`
	DonatedInternalID2 string `db:"-" csv:"DONATED_INTERNAL_ID2"`
	CardNumber         string `db:"-" csv:"CARD_NUMBER"`
	PayParams          string `db:"-" csv:"PAY_PARAMS"`
	PersonRecieved     string `db:"-" csv:"PERSON_RECIEVED"`
	BankDivisionName   string `db:"-" csv:"BANK_DIVISION_NAME"`
	BankCardID         string `db:"-" csv:"BANK_CARD_ID"`
	AddressTypeID      string `db:"-" csv:"ADDRESS_TYPE_ID"`
	AddressType        string `db:"-" csv:"ADDRESS_TYPE"`
	Zip                string `db:"-" csv:"ZIP"`
	Country            string `db:"-" csv:"COUNTRY"`
	Region             string `db:"-" csv:"REGION"`
	Zone               string `db:"-" csv:"ZONE"`
	City               string `db:"-" csv:"CITY"`
	Street             string `db:"-" csv:"STREET"`
	Building           string `db:"-" csv:"BUILDING"`
	BuildSect          string `db:"-" csv:"BUILD_SECT"`
	Apartment          string `db:"-" csv:"APARTMENT"`
	UnstructInfo       string `db:"-" csv:"UNSTRUCT_INFO"`
}

// Payments is a generator for PAYMENTS table
type Payments struct{}
