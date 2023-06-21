package gen

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/csv"
)

type IPPlanRow struct {
	Description string    `db:"name" csv:"DESCRIPTION"`
	IPType      int       `db:"-" csv:"IP_TYPE"`
	IPv4        string    `db:"network" csv:"IPV4"`
	IPv6        string    `db:"-" csv:"IPV6"`
	IPv4Mask    string    `db:"mask" csv:"IPV4_MASK"`
	IPv6Mask    string    `db:"-" csv:"IPV6_MASK"`
	BeginTime   time.Time `db:"-" csv:"BEGIN_TIME" time:"2006-01-02 15:04:05"`
	EndTime     time.Time `db:"-" csv:"END_TIME" time:"2006-01-02 15:04:05"`
	RegionID    int       `db:"-" csv:"REGION_ID"`
}

type IPPlan struct{}

func (a *IPPlan) Render(db *sqlx.DB) (r []string, err error) {
	var plan []IPPlanRow
	err = db.Select(&plan, `select INET_NTOA(network) as network,
		INET_NTOA(mask) as mask,
		name from dhcphosts_networks order by id`)
	if err != nil {
		return nil, err
	}
	for i := range plan {
		plan[i].Calc()
	}
	r = csv.MarshalCSV(plan, ";", "")
	return r, nil
}

func (a *IPPlan) GetFileName() string {
	return fmt.Sprintf("IP_PLAN_%s.txt", time.Now().Format("20060102_1504"))
}

func (a *IPPlanRow) Calc() {
	a.IPType = 0
	a.RegionID = EnvRegionID
	a.BeginTime = EnvInitDate
	if a.IPv4 != "" && a.IPv4Mask != "" {
		a.IPv4 = MakeIP(a.IPv4)
		a.IPv4Mask = MakeIP(a.IPv4Mask)
	}
}
