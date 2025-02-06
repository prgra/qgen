package gen

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/config"
	"github.com/prgra/qgen/csv"
)

type AbonUsers struct {
}

type AbonUsersRow struct {
	AbonID       int    `db:"-" csv:"ABONENT_ID"`
	RegionID     int    `csv:"REGION_ID"`
	UserNumber   string `db:"phone" csv:"USER_NUMBER"`
	UserName     string `db:"fio" csv:"USER_NAME"`
	RecordAction int    `csv:"RECORD_ACTION"`
	InternalID1  string `db:"uid" csv:"INTERNAL_ID1"`
	Company      int    `db:"company_id" csv:"-"`
}

func (a *AbonUsers) Render(db *sqlx.DB, cfg config.Config) (r []string, err error) { //
	var abons []AbonUsersRow //
	dta := cfg.InitDate.Format("2006-01-02")
	if cfg.OnlyOneDay {
		dta = time.Now().Format("2006-01-02")
	}
	err = db.Select(&abons, `SELECT u.company_id, u.uid, pi.fio as fio, pi.phone as phone FROM companies c
	JOIN users u ON u.company_id = c.id
	JOIN users_pi pi on u.uid = pi.uid and c.registration >= ?`, dta)
	if err != nil {
		return nil, err
	}
	for i := range abons {
		abons[i].Calc(cfg)
	}
	r = csv.MarshalCSV(abons, ";", "")
	return r, nil
}

func (a *AbonUsers) GetFileName() string {
	return fmt.Sprintf("ABONENT_USER_%s.txt", time.Now().Format("20060102_1504"))
}

func (a *AbonUsers) GetRemoteDir() string {
	return ""
}

func (r *AbonUsersRow) Calc(cfg config.Config) {
	r.InternalID1 = fmt.Sprintf("%s%d", cfg.CompanyCode, r.Company)
	r.RegionID = cfg.RegionID
	r.RecordAction = 1
}
