package gen

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/csv"
)

type AbonAddrRow struct {
	AbonentID     int            `db:"uid" csv:"ABONENT_ID"`
	RegionID      sql.NullInt64  `db:"district_id" csv:"REGION_ID"`
	AddressTypeID int            `db:"-" csv:"ADDRESS_TYPE_ID"`
	AddressType   int            `db:"-" csv:"ADDRESS_TYPE"`
	Zip           string         `db:"-" csv:"ZIP"`
	Country       string         `db:"-" csv:"COUNTRY"`
	Region        string         `db:"-" csv:"REGION"`
	Zone          string         `db:"-" csv:"ZONE"`
	City          sql.NullString `db:"dist" csv:"CITY"`
	Street        sql.NullString `db:"street" csv:"STREET"`
	Building      sql.NullString `db:"build" csv:"BUILDING"`
	BuildSect     string         `db:"-" csv:"BUILD_SECT"`
	Apartment     string         `db:"flat" csv:"APARTMENT"`
	UnstructInfo  string         `db:"-" csv:"UNSTRUCT_INFO"`
	BeginTime     string         `db:"-" csv:"BEGIN_TIME"`
	EndTime       string         `db:"-" csv:"END_TIME"`
	RecordAction  string         `db:"-" csv:"RECORD_ACTION"`
	InternalID1   string         `db:"id" csv:"INTERNAL_ID1"`
	InternalID2   string         `db:"-" csv:"INTERNAL_ID2"`
}

type AbonAddr struct{}

func (a *AbonAddr) Render(db *sqlx.DB) (r []string, err error) { //
	var abons []AbonAddrRow //
	err = db.Select(&abons, `select u.uid, 
s.district_id,
u.id,
d.name as dist,
s.name as street,
b.number as build,
pi.address_flat as flat

from 
users u 
JOIN dv_main dv ON dv.uid=u.uid
LEFT JOIN users_pi pi ON pi.uid=u.uid 
LEFT JOIN builds b ON b.id=pi.location_id
LEFT JOIN streets s ON s.id=b.street_id
LEFT JOIN districts d ON d.id=s.district_id
LEFT JOIN bills bi ON u.bill_id=bi.id
LEFT JOIN companies c ON c.id=u.company_id
JOIN tarif_plans tp ON tp.id=dv.tp_id
`)
	if err != nil {
		return nil, err
	}
	for i := range abons {
		abons[i].Calc()
	}
	r = csv.MarshalCSV(abons, ";", "")
	return r, nil
}

func (a *AbonAddr) GetFileName() string {
	return fmt.Sprintf("ABONENT_ADDR_%s.txt", time.Now().Format("20060102_1504"))
}

func (a *AbonAddrRow) Calc() {
	a.AddressTypeID = 0
	a.AddressType = 1
	if a.City.Valid && a.Street.Valid {
		a.UnstructInfo = fmt.Sprintf("%s %s %s %s", a.City.String, a.Street.String, a.Building.String, a.Apartment)
	}
}
