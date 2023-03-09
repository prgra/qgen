package csv

import (
	"testing"
	"time"
)

func TestCSV(t *testing.T) {
	type C struct {
		ID    int `csv:"-"`
		Name  string
		Age   int       `csv:"возраст"`
		Birth time.Time `csv:"birth" time:"2006-01-02"`
	}
	var csv = []C{
		{ID: 1, Name: "John", Age: 30, Birth: time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)},
		{ID: 2, Name: "Mary", Age: 25, Birth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)},
		{ID: 3, Name: "Mike", Age: 35, Birth: time.Date(1995, 1, 1, 0, 0, 0, 0, time.UTC)},
	}
	MarshalCSV(csv, ";", "'")
	t.Fail()
}
