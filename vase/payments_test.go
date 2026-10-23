package vase

import (
	"testing"
	"time"

	"github.com/prgra/qgen/csv"
)

func TestPaymentRow(t *testing.T) {
	msk := time.FixedZone("MSK", 3*60*60)
	p := PaymentRow{
		ID:          "123",
		Login:       "login_123",
		PaymentDate: time.Date(2024, 3, 1, 12, 30, 0, 0, msk),
		Amount:      1500.5,
		Comment:     "оплата\r\nчерез кассу",
	}
	p.Calc()
	got := csv.MarshalCSVNoHeader([]PaymentRow{p}, ";", "\"")
	want := `"123";"login_123";"01.03.2024T09:30:00";"1500,50";"оплата через кассу"`
	if len(got) != 1 || got[0] != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
