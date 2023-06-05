package gen

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/csv"
)

type GateWayRow struct {
	GateID        int       `csv:"GATE_ID"`
	BeginTime     time.Time `csv:"BEGIN_TIME" time:"2006-01-02 15:04:05"`
	EndTime       time.Time `csv:"END_TIME" time:"2006-01-02 15:04:05"`
	Description   string    `csv:"DESCRIPTION"`
	GateType      int       `csv:"GATE_TYPE"`
	AddressTypeID int       `csv:"ADDRESS_TYPE_ID"`
	AddressType   int       `csv:"ADDRESS_TYPE"`
	Zip           string    `csv:"ZIP"`
	Country       string    `csv:"COUNTRY"`
	Region        string    `csv:"REGION"`
	Zone          string    `csv:"ZONE"`
	City          string    `csv:"CITY"`
	Street        string    `csv:"STREET"`
	Building      string    `csv:"BUILDING"`
	BuildSect     string    `csv:"BUILD_SECT"`
	Apartment     string    `csv:"APARTMENT"`
	UnstructInfo  string    `csv:"UNSTRUCT_INFO"`
	RegionID      int       `csv:"REGION_ID"`
}

type GateWay struct{}

func (a *GateWay) Render(db *sqlx.DB) (r []string, err error) {
	gt := []GateWayRow{
		{
			GateID:      1,
			BeginTime:   EnvInitDate,
			Description: "NAS",
			RegionID:    EnvRegionID,
		},
	}

	r = csv.MarshalCSV(gt, ";", "")
	return r, nil
}

func (a *GateWay) GetFileName() string {
	return fmt.Sprintf("GATEWAY_%s.txt", time.Now().Format("20060102_1504"))
}
