package gen

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/csv"
)

type RegionRow struct {
	ID          string `db:"id" csv:"ID"`
	BEGIN_TIME  string `db:"-" csv:"BEGIN_TIME"`
	END_TIME    string `db:"-" csv:"END_TIME"`
	DESCRIPTION string `db:"name" csv:"DESCRIPTION"`
	MCC         string `db:"-" csv:"MCC"`
	MNC         string `db:"-" csv:"MNC"`
}

type Region struct{}

func (a *Region) Render(db *sqlx.DB) (r []string, err error) {
	var regions []RegionRow
	err = db.Select(&regions, `select id, name from districts`)
	if err != nil {
		return nil, err
	}
	r = csv.MarshalCSV(regions, ";", "\"")
	return r, nil
}

func (a *Region) GetFileName() string {
	return fmt.Sprintf("REGION_%s.txt", time.Now().Format("20060102_1504"))
}
