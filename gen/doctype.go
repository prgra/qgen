package gen

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/config"
	"github.com/prgra/qgen/csv"
)

type DocTypeRow struct {
	DocTypeID   int       `csv:"DOC_TYPE_ID"`
	BeginTime   time.Time `csv:"BEGIN_TIME" time:"2006-01-02 15:04:05"`
	EndTime     time.Time `csv:"END_TIME" time:"2006-01-02 15:04:05"`
	Description string    `csv:"DESCRIPTION"`
	RegionID    int       `csv:"REGION_ID"`
}

type DocType struct {
}

func (a *DocType) Render(db *sqlx.DB, cfg config.Config) (r []string, err error) {
	regions := []DocTypeRow{
		{
			DocTypeID:   1,
			BeginTime:   cfg.InitDate,
			Description: "паспорт",
			RegionID:    cfg.RegionID,
		},
		{
			DocTypeID:   2,
			BeginTime:   cfg.InitDate,
			Description: "другое",
			RegionID:    cfg.RegionID,
		},
	}

	r = csv.MarshalCSV(regions, ";", "")
	return r, nil
}

func (a *DocType) GetFileName() string {
	return fmt.Sprintf("DOC_TYPE_%s.txt", time.Now().Format("20060102_1504"))
}

func (a *DocType) GetRemoteDir() string {
	return ""
}
