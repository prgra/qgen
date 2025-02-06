package yhnt

import (
	"encoding/binary"
	"net"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/config"
	"github.com/prgra/qgen/csv"
)

// IPPlanRow is a row from the plan table
type IPPlanRow struct {
	DepID     int       ``                           // 1 идентификатор филиала (число);
	Descr     string    `db:"name"`                  // 2 описание назначения диапазона IP-адресов (строка);
	Net       string    `db:"net"`                   // 3 IP-адрес подсети (либо v4 либо v6);
	Mask      int       `db:"mask"`                  // 4 маска подсети (число);
	StartDate time.Time `time:"02.01.2006 15:04:05"` // 5 время начала действия диапазона (дата);
	EndTime   time.Time `time:"02.01.2006 15:04:05"` // 6 дата завершения действия диапазона (в случае работы оборудования по настоящее время - пустое значение).
}

// IPPlan is a generator for ipplan table
type IPPlan struct{}

func (a *IPPlan) Render(db *sqlx.DB, cfg config.Config) (r []string, err error) {
	var ipplan []IPPlanRow

	err = db.Select(&ipplan, `SELECT name, INET_NTOA(network) as net,mask FROM dhcphosts_networks`)
	if err != nil {
		return nil, err
	}
	for i := range ipplan {
		ipplan[i].Calc(cfg)
	}
	r = csv.MarshalCSVNoHeader(ipplan, ";", "\"")
	return r, nil
}

func (p *IPPlanRow) Calc(cfg config.Config) {
	if p.Mask > 32 {
		a := make([]byte, 4)
		binary.BigEndian.PutUint32(a, 4294901760)
		ms := net.IPv4Mask(a[0], a[1], a[2], a[3])
		p.Mask, _ = ms.Size()
	}
	p.StartDate = cfg.InitDate
	p.DepID = cfg.RegionID
}

func (a *IPPlan) GetFileName() string {
	return "ip-numbering-plan.csv"
}

func (a *IPPlan) GetRemoteDir() string {
	return ""
}
