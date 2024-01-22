package gen

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/config"
	"github.com/prgra/qgen/csv"
)

type RegionRow struct {
	ID          int       `db:"id" csv:"ID"`
	BeginTime   time.Time `db:"-" csv:"BEGIN_TIME" time:"2006-01-02 15:04:05"`
	EndTime     time.Time `db:"-" csv:"END_TIME" time:"2006-01-02 15:04:05"`
	DESCRIPTION string    `db:"name" csv:"DESCRIPTION"`
	MCC         string    `db:"-" csv:"MCC"`
	MNC         string    `db:"-" csv:"MNC"`
}

type Region struct{}

func (a *Region) RenderMysql(db *sqlx.DB, cfg config.Config) (r []string, err error) {
	var regions []RegionRow
	err = db.Select(&regions, `select id, name from districts order by id`)
	if err != nil {
		return nil, err
	}
	var rr []RegionRow
	rr = append(rr, RegionRow{
		ID:          0,
		BeginTime:   cfg.InitDate,
		DESCRIPTION: "Не указан",
	})
	for i := range regions {
		regions[i].BeginTime = cfg.InitDate
	}
	rr = append(rr, regions...)
	r = csv.MarshalCSV(rr, ";", "")
	return r, nil
}

func (a *Region) Render(db *sqlx.DB, cfg config.Config) (r []string, err error) {

	r = csv.MarshalCSV([]RegionRow{
		{
			ID:          cfg.RegionID,
			BeginTime:   cfg.InitDate,
			DESCRIPTION: cfg.RegionName,
		},
	}, ";", "")
	return r, nil
}

func (a *Region) GetFileName() string {
	return fmt.Sprintf("REGIONS_%s.txt", time.Now().Format("20060102_1504"))
}
