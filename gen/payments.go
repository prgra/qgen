package gen

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/csv"
)

// PaymentRow is a row from the PAYMENTS table
type PaymentRow struct {
	RegionID           int            `db:"-" csv:"REGION_ID"`
	PaymentType        int            `db:"-" csv:"PAYMENT_TYPE"`
	PayTypeID          int            `db:"method" csv:"PAY_TYPE_ID"`
	PaymentDate        time.Time      `db:"date" csv:"PAYMENT_DATE" time:"2006-01-02 15:04:05"`
	Amount             string         `db:"sum" csv:"AMOUNT"`
	AmountCurrency     string         `db:"-" csv:"AMOUNT_CURRENCY"`
	PhoneNumber        sql.NullString `db:"phone" csv:"PHONE_NUMBER"`
	Account            string         `db:"bill_id" csv:"ACCOUNT"`
	AbonentID          string         `db:"uid" csv:"ABONENT_ID"`
	BankAccount        string         `db:"-" csv:"BANK_ACCOUNT"`
	BankName           string         `db:"-" csv:"BANK_NAME"`
	ExpressCardNumber  string         `db:"-" csv:"EXPRESS_CARD_NUMBER"`
	TerminalID         string         `db:"-" csv:"TERMINAL_ID"`
	TerminalNumber     string         `db:"-" csv:"TERMINAL_NUMBER"`
	LATITUDE           string         `db:"-" csv:"LATITUDE"`
	LONGITUDE          string         `db:"-" csv:"LONGITUDE"`
	ProjectionType     string         `db:"-" csv:"PROJECTION_TYPE"`
	CenterID           string         `db:"-" csv:"CENTER_ID"`
	DonatedPhoneNumber string         `db:"-" csv:"DONATED_PHONE_NUMBER"`
	DonatedAccount     string         `db:"-" csv:"DONATED_ACCOUNT"`
	DonatedInternalID1 string         `db:"-" csv:"DONATED_INTERNAL_ID1"`
	DonatedInternalID2 string         `db:"-" csv:"DONATED_INTERNAL_ID2"`
	CardNumber         string         `db:"-" csv:"CARD_NUMBER"`
	PayParams          string         `db:"-" csv:"PAY_PARAMS"`
	PersonReceived     string         `db:"-" csv:"PERSON_RECIEVED"`
	BankDivisionName   string         `db:"-" csv:"BANK_DIVISION_NAME"`
	BankCardID         string         `db:"-" csv:"BANK_CARD_ID"`
	AddressTypeID      int            `db:"-" csv:"ADDRESS_TYPE_ID"`
	AddressType        string         `db:"-" csv:"ADDRESS_TYPE"`
	Zip                string         `db:"-" csv:"ZIP"`
	Country            string         `db:"-" csv:"COUNTRY"`
	Region             string         `db:"-" csv:"REGION"`
	Zone               string         `db:"-" csv:"ZONE"`
	City               string         `db:"-" csv:"CITY"`
	Street             string         `db:"-" csv:"STREET"`
	Building           string         `db:"-" csv:"BUILDING"`
	BuildSect          string         `db:"-" csv:"BUILD_SECT"`
	Apartment          string         `db:"-" csv:"APARTMENT"`
	UnstructInfo       string         `db:"-" csv:"UNSTRUCT_INFO"`
	RecordAction       int            `db:"-" csv:"RECORD_ACTION"`
}

// Payments is a generator for PAYMENTS table
type Payments struct{}

func (a *Payments) Render(db *sqlx.DB) (r []string, err error) {
	var pays []PaymentRow
	dta := EnvInitDate.Format("2006-01-02")
	if EnvOnlyOneDay {
		dta = time.Now().Format("2006-01-02")
	}
	pmap, _, err := LoadPayMethodsMapFromFile("paymethods.map")
	if err != nil {
		return nil, err
	}
	err = db.Select(&pays, `SELECT p.date, p.sum, p.method, p.uid, p.bill_id, pi.phone 
		FROM payments p
		LEFT JOIN users_pi pi on pi.uid = p.uid
		where date >= ? order by id`, dta)
	if err != nil {
		return nil, err
	}
	for i := range pays {
		pays[i].Calc(pmap)
	}
	r = csv.MarshalCSV(pays, ";", "")
	return r, nil
}

func (p *PaymentRow) Calc(pmap map[int]int) {
	p.RegionID = EnvRegionID
	p.Country = EnvCountry
	p.PaymentType = pmap[p.PayTypeID]
	p.AmountCurrency = p.Amount
	p.RecordAction = 1

}

// LoadMethodsMapFromFile loads map of payment methods from file
// file struct is:
// 0:83:Наличные
// 1:80:Банк
// 2:82:Внешние платежи
// 3:80:Credit Card
// 4:86:Бонус
// 5:86:Корректировка
func LoadPayMethodsMapFromFile(filename string) (r map[int]int, pt []PayTypeRow, err error) {
	f, err := os.Open(filepath.Clean(filename))
	if err != nil {
		return nil, nil, err
	}
	fa, err := io.ReadAll(f)
	if err != nil {
		return nil, nil, err
	}
	r = make(map[int]int)
	for _, v := range strings.Split(string(fa), "\n") {
		a := strings.Split(v, ":")
		if len(a) != 3 {
			continue
		}
		i1, err1 := strconv.Atoi(a[0])
		i2, err2 := strconv.Atoi(a[1])
		if err1 != nil || err2 != nil {
			continue
		}
		r[i1] = i2
		var p PayTypeRow
		p.ID = i1
		p.Descr = fmt.Sprintf(a[2])
		p.BeginTime = EnvInitDate
		p.RegionID = EnvRegionID
		pt = append(pt, p)
	}
	return r, pt, nil
}

func (a *Payments) GetFileName() string {
	return fmt.Sprintf("PAYMENT_%s.txt", time.Now().Format("20060102_1504"))
}
