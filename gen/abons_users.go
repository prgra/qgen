package gen

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
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
	InternalID1  string `db:"id" csv:"INTERNAL_ID1"`
}

func (a *AbonUsers) Render(db *sqlx.DB) (r []string, err error) { //
	var abons []AbonUsersRow //
	dta := EnvInitDate.Format("2006-01-02")
	if EnvOnlyOneDay {
		dta = time.Now().Format("2006-01-02")
	}
	err = db.Select(&abons, `SELECT c.id as id, pi.fio as fio, pi.phone as phone FROM companies c
	JOIN users u ON u.company_id = c.id
	JOIN users_pi pi on u.uid = pi.uid and c.registration >= ?`, dta)
	if err != nil {
		return nil, err
	}
	for i := range abons {
		abons[i].Calc()
	}
	r = csv.MarshalCSV(abons, ";", "")
	return r, nil
}

func (a *AbonUsers) GetFileName() string {
	return fmt.Sprintf("ABONENT_USER_%s.txt", time.Now().Format("20060102_1504"))
}

func (r *AbonUsersRow) Calc() {
	r.RegionID = EnvRegionID
}
