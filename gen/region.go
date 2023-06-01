package gen

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/csv"
)

type RegionRow struct {
	ID          string    `db:"id" csv:"ID"`
	BeginTime   time.Time `db:"-" csv:"BEGIN_TIME" time:"2006-01-02 15:06:07"`
	EndTime     time.Time `db:"-" csv:"END_TIME" time:"2006-01-02 15:06:07"`
	DESCRIPTION string    `db:"name" csv:"DESCRIPTION"`
	MCC         string    `db:"-" csv:"MCC"`
	MNC         string    `db:"-" csv:"MNC"`
}

type Region struct{}

func (a *Region) Render(db *sqlx.DB) (r []string, err error) {
	var regions []RegionRow
	err = db.Select(&regions, `select id, name from districts order by id`)
	if err != nil {
		return nil, err
	}
	var rr []RegionRow
	rr = append(rr, RegionRow{
		ID:          "0",
		BeginTime:   time.Unix(0, 0),
		DESCRIPTION: "Не указан",
	})
	for i := range regions {
		regions[i].BeginTime = time.Unix(0, 0)
	}
	rr = append(rr, regions...)
	r = csv.MarshalCSV(rr, ";", "")
	return r, nil
}

func (a *Region) GetFileName() string {
	return fmt.Sprintf("REGIONS_%s.txt", time.Now().Format("20060102_1504"))
}
