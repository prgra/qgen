package gen

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/csv"
)

type DocTypeRow struct {
	DocTypeID   int       `csv:"DOC_TYPE_ID"`
	BeginTime   time.Time `csv:"BEGIN_TIME"`
	EndTime     time.Time `csv:"END_TIME"`
	Description string    `csv:"DESCRIPTION"`
	RegionID    int       `csv:"REGION_ID"`
}

type DocType struct{}

func (a *DocType) Render(db *sqlx.DB) (r []string, err error) {
	regions := []DocTypeRow{
		{
			DocTypeID:   1,
			BeginTime:   time.Unix(0, 0),
			Description: "паспорт",
		},
	}

	r = csv.MarshalCSV(regions, ";", "")
	return r, nil
}

func (a *DocType) GetFileName() string {
	return fmt.Sprintf("SUPPLEMENTARY_SERVICE_%s.txt", time.Now().Format("20060102_1504"))
}
