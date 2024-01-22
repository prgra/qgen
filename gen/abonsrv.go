package gen

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/config"
	"github.com/prgra/qgen/csv"
)

type AbonSrv struct{}

type AbonSrvRow struct {
	AbonentID    sql.NullInt64 `db:"" csv:"ABONENT_ID"`
	RegionID     int           `db:"" csv:"REGION_ID"`
	ID           int           `db:"" csv:"ID"`
	BeginTime    time.Time     `db:"attach" csv:"BEGIN_TIME"`
	EndTime      time.Time     `db:"" csv:"END_TIME"`
	Parameter    string        `db:"" csv:"PARAMETER"`
	SrvContract  string        `db:"login" csv:"SRV_CONTRACT"`
	RecordAction int           `db:"" csv:"RECORD_ACTION"`
	InternalID1  string        `db:"uid" csv:"INTERNAL_ID1"`
	InternalID2  string        `db:"" csv:"INTERNAL_ID2"`
	Company      int           `db:"company_id" csv:"-"`
}

func (a *AbonSrv) Render(db *sqlx.DB, cfg config.Config) (r []string, err error) { //
	var abons []AbonSrvRow //
	dta := cfg.InitDate.Format("2006-01-02")
	if cfg.OnlyOneDay {
		dta = time.Now().Format("2006-01-02")
	}
	err = db.Select(&abons, `SELECT
u.id as login,
u.uid, 
aa1.datetime as attach,
u.company_id
FROM 
users u 
JOIN dv_main dv ON dv.uid=u.uid
LEFT JOIN admin_actions aa1 on aa1.id = (select id from admin_actions 
	where uid=u.uid order by id desc limit 1)
WHERE aa1.datetime >= ?`, dta)
	if err != nil {
		return nil, err
	}
	for i := range abons {
		abons[i].Calc(cfg)
	}
	r = csv.MarshalCSV(abons, ";", "")
	return r, nil
}

func (a *AbonSrv) GetFileName() string {
	return fmt.Sprintf("ABONENT_SRV_%s.txt", time.Now().Format("20060102_1504"))
}

func (a *AbonSrvRow) Calc(cfg config.Config) {
	if a.Company > 0 {
		a.InternalID1 = fmt.Sprintf("%s%d", cfg.CompanyCode, a.Company)
	}
	a.RegionID = cfg.RegionID
	a.RecordAction = 1
	a.ID = 1
}
