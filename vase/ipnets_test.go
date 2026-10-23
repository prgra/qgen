package vase

import (
	"database/sql"
	"testing"
	"time"

	"github.com/prgra/qgen/csv"
)

func TestIPNetRow(t *testing.T) {
	begin := sql.NullTime{Time: time.Date(2017, 1, 1, 1, 1, 1, 0, time.UTC), Valid: true}
	end := sql.NullTime{Time: time.Date(2020, 2, 2, 2, 2, 2, 0, time.UTC), Valid: true}
	tests := []struct {
		name string
		row  IPNetRow
		want string
	}{
		{name: "серый адрес с MAC", row: IPNetRow{ID: "1", IP: "192.168.0.1", Mask: sql.NullString{String: "255.255.255.255", Valid: true},
			MAC: sql.NullString{String: "0a-0b-0c-0d-0e-0f", Valid: true}, ContractDate: begin},
			want: `"";"1";"192.168.0.1";"0";"01.01.2017T01:01:01";"";"0A:0B:0C:0D:0E:0F";""`},
		{name: "публичная подсеть без MAC, отключен", row: IPNetRow{ID: "2", IP: "85.202.1.1", Mask: sql.NullString{String: "255.255.255.252", Valid: true},
			MAC: sql.NullString{String: "0a:0b:0c:0d:0e:0f", Valid: true}, Status: 1, FirstAction: begin, LastAction: end},
			want: `"";"2";"85.202.1.0/30";"1";"01.01.2017T01:01:01";"02.02.2020T02:02:02";"";""`},
		{name: "CGNAT без маски", row: IPNetRow{ID: "3", IP: "100.64.1.10"},
			want: `"";"3";"100.64.1.10";"0";"";"";"";""`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.row.Calc()
			got := csv.MarshalCSVNoHeader([]IPNetRow{tt.row}, ";", "\"")
			if len(got) != 1 || got[0] != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
