package vase

import (
	"database/sql"
	"testing"
	"time"

	"github.com/prgra/qgen/config"
	"github.com/prgra/qgen/csv"
)

func TestServiceList(t *testing.T) {
	cfg := config.Config{InitDate: time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC)}
	got, err := (&ServiceList{}).Render(nil, cfg)
	want := `"1";"Интернет";"Высокоскоростной доступ в интернет";"01.01.2016";"";"Ethernet"`
	if err != nil || len(got) != 1 || got[0] != want {
		t.Errorf("got %q, %v, want %q", got, err, want)
	}
}

func TestServiceRow(t *testing.T) {
	begin := sql.NullTime{Time: time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC), Valid: true}
	end := sql.NullTime{Time: time.Date(2019, 3, 4, 5, 6, 7, 0, time.UTC), Valid: true}
	tests := []struct {
		name string
		row  ServiceRow
		want string
	}{
		{name: "активный", row: ServiceRow{ID: "301", ContractDate: begin, LastAction: end},
			want: `"1";"301";"";"";"";"Высокоскоростной доступ в интернет";"01.01.2017T00:00:00";"";""`},
		{name: "отключен", row: ServiceRow{ID: "302", Status: 1, ContractDate: begin, LastAction: end},
			want: `"1";"302";"";"";"";"Высокоскоростной доступ в интернет";"01.01.2017T00:00:00";"04.03.2019T05:06:07";""`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.row.Calc()
			got := csv.MarshalCSVNoHeader([]ServiceRow{tt.row}, ";", "\"")
			if len(got) != 1 || got[0] != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
