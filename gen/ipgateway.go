package gen

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/csv"
)

type IPGateWayRow struct {
	GateID   int    `db:"id" csv:"GATE_ID"`
	IPType   int    `csv:"IP_TYPE"`
	IPv4     string `db:"ip" csv:"IPV4"`
	IPv6     string `csv:"IPV6"`
	IPPort   string `csv:"IP_PORT"`
	RegionID int    `csv:"REGION_ID"`
}

type IPGateWay struct{}

func (a *IPGateWay) Render(db *sqlx.DB, cfg Config) (r []string, err error) {
	var gws []IPGateWayRow
	err = db.Select(&gws, `select id,ip from nas where gid = ? order by id`, cfg.NasGroupID)
	if err != nil {
		return nil, err
	}
	for i := range gws {
		gws[i].Calc(cfg)
	}
	r = csv.MarshalCSV(gws, ";", "")
	return r, nil

	// gt := []GateWayRow{
	// 	{
	// 		GateID:      1,
	// 		BeginTime:   cfg.InitDate,
	// 		GateType:    8,
	// 		Description: "NAS",
	// 		RegionID:    cfg.RegionID,
	// 	},
	// }

	// r = csv.MarshalCSV(gt, ";", "")
	// return r, nil
}

func (a *IPGateWayRow) Calc(cfg Config) {
	a.RegionID = cfg.RegionID
	a.IPType = 0
	a.IPv4 = MakeIP(a.IPv4)
}

func (a *IPGateWay) GetFileName() string {
	return fmt.Sprintf("IP_GATEWAY_%s.txt", time.Now().Format("20060102_1504"))
}
