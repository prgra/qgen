package yhnt

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/config"
	"github.com/prgra/qgen/csv"
)

type DocTypeRow struct {
	RegionID    int       `csv:"REGION_ID"`
	DocTypeID   int       `csv:"DOC_TYPE_ID"`
	BeginTime   time.Time `csv:"BEGIN_TIME" time:"02.01.2006 15:04:05"`
	EndTime     time.Time `csv:"END_TIME" time:"02.01.2006 15:04:05"`
	Description string    `csv:"DESCRIPTION"`
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

	r = csv.MarshalCSVNoHeader(regions, ";", "\"")
	return r, nil
}

func (a *DocType) GetFileName() string {
	return "doc_types.csv"
}

func (a *DocType) GetRemoteDir() string {
	return ""
}
