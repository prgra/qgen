package vase

import (
	"testing"
	"time"

	"github.com/prgra/qgen/csv"
)

func TestPaymentRow(t *testing.T) {
	methods := map[int]string{0: "Наличные", 67: "Sberbank"}
	tests := []struct {
		name string
		row  PaymentRow
		want string
	}{
		{name: "безналичный", row: PaymentRow{ID: "301", PayID: 15, Method: 67, Amount: 500,
			PaymentDate: time.Date(2018, 1, 1, 18, 15, 15, 0, time.FixedZone("MSK", 3*60*60))},
			want: `"301";"";"01.01.2018T15:15:15";"500,00";"Безналичный";"Sberbank, ID платежа 15";""`},
		{name: "наличный с комментарием", row: PaymentRow{ID: "302", PayID: 16, Method: 0, Amount: 1500.5,
			Comment: "касса\r\nофис 1", PaymentDate: time.Date(2018, 1, 1, 15, 15, 15, 0, time.UTC)},
			want: `"302";"";"01.01.2018T15:15:15";"1500,50";"Наличный";"Наличные, ID платежа 16, касса офис 1";""`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.row.Calc(methods)
			got := csv.MarshalCSVNoHeader([]PaymentRow{tt.row}, ";", "\"")
			if len(got) != 1 || got[0] != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
