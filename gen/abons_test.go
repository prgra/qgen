package gen

import (
	"os"
	"testing"

	"github.com/prgra/qgen/csv"
)

func TestAbons(t *testing.T) {
	var ab []AbonRow
	for i := 0; i < 10; i++ {
		a := AbonRow{
			AbonID: i,
			FIO:    "ываыв",
		}
		a.Calc()
		ab = append(ab, a)
	}
	str := csv.MarshalCSV(ab, ",", "'")
	f, _ := os.Create("abons.csv")
	for _, s := range str {
		f.Write([]byte(s + "\n"))
	}
	f.Close()

	t.Fail()

}
