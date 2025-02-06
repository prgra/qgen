package gen

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/config"
	"github.com/prgra/qgen/csv"
)

type PayTypes struct{}

type PayTypeRow struct {
	ID        int       `csv:"ID"`
	BeginTime time.Time `csv:"BEGIN_TIME" format:"2006-01-02 15:04:05"`
	EndTime   time.Time `csv:"END_TIME"`
	Descr     string    `csv:"DESCRIPTION"`
	RegionID  int       `csv:"REGION_ID"`
}

func (a *PayTypes) Render(db *sqlx.DB, cfg config.Config) (r []string, err error) {
	_, tps, err := LoadPayMethodsMapFromFile("paymethods.map", cfg)
	if err != nil {
		return nil, err
	}
	r = csv.MarshalCSV(tps, ";", "")
	return r, nil
}

func (a *PayTypes) GetFileName() string {
	return fmt.Sprintf("PAY_TYPE_%s.txt", time.Now().Format("20060102_1504"))
}

func (a *PayTypes) GetRemoteDir() string {
	return ""
}
