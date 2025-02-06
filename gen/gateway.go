package gen

import (
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/config"
	"github.com/prgra/qgen/csv"
)

type GateWayRow struct {
	GateID        int       `db:"id" csv:"GATE_ID"`
	BeginTime     time.Time `csv:"BEGIN_TIME" time:"2006-01-02 15:04:05"`
	EndTime       time.Time `csv:"END_TIME" time:"2006-01-02 15:04:05"`
	Description   string    `db:"name" csv:"DESCRIPTION"`
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
	UnstructInfo  string    `db:"addr" csv:"UNSTRUCT_INFO"`
	RegionID      int       `csv:"REGION_ID"`
}

type GateWay struct{}

func (a *GateWay) Render(db *sqlx.DB, cfg config.Config) (r []string, err error) {

	var gws []GateWayRow
	err = db.Select(&gws, `select n.id, n.name, 
	CONCAT(d.name, ' ', s.name, ' ', b.number, ' ', n.address_flat) as addr
	from nas n
	left join builds b on b.id = n.location_id
	left join streets s on s.id = b.street_id
	left join districts d on d.id = s.district_id
	where n.gid = ? order by n.id`, cfg.NasGroupID)
	if err != nil {
		return nil, err
	}
	for i := range gws {
		gws[i].Calc(cfg)
	}
	r = csv.MarshalCSV(gws, ";", "")
	return r, nil
}

func (a *GateWayRow) Calc(cfg config.Config) {
	a.AddressTypeID = 3
	a.AddressType = 1
	a.GateType = 8
	a.RegionID = cfg.RegionID
	a.BeginTime = cfg.InitDate
	a.UnstructInfo = strings.TrimSpace(a.UnstructInfo)
}

func (a *GateWay) GetFileName() string {
	return fmt.Sprintf("GATEWAY_%s.txt", time.Now().Format("20060102_1504"))
}

func (a *GateWay) GetRemoteDir() string {
	return ""
}
