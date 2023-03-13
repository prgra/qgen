package gen

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/csv"
)

//
//
//
//
//

type AbonIdent struct{}

type AbonIdentRow struct {
	AbonID            int            `db:"uid" csv:"ABONENT_ID"`
	RegionID          sql.NullInt64  `db:"district_id" csv:"REGION_ID"`
	IdentType         int            `db:"-" csv:"IDENT_TYPE"`
	Phone             string         `db:"-" csv:"PHONE"`
	InternalNumber    string         `db:"-" csv:"INTERNAL_NUMBER"`
	IMSI              string         `db:"-" csv:"IMSI"`
	IMEI              string         `db:"-" csv:"IMEI"`
	ICC               string         `db:"-" csv:"ICC"`
	MIN               string         `db:"-" csv:"MIN"`
	ESN               string         `db:"-" csv:"ESN"`
	EquipmentType     int            `db:"-" csv:"EQUIPMENT_TYPE"`
	MAC               sql.NullString `db:"mac" csv:"MAC"`
	VPI               string         `db:"-" csv:"VPI"`
	VCI               string         `db:"-" csv:"VCI"`
	Login             string         `db:"id" csv:"LOGIN"`
	EMail             string         `db:"email" csv:"E_MAIL"`
	PIN               string         `db:"-" csv:"PIN"`
	UserDomain        string         `db:"-" csv:"USER_DOMAIN"`
	Reserved          string         `db:"-" csv:"RESERVED"`
	OriginatorName    string         `db:"-" csv:"ORIGINATOR_NAME"`
	IPType            int            `db:"-" csv:"IP_TYPE"`
	IPv4              string         `db:"ip" csv:"IPV4"`
	IPv6              string         `db:"-" csv:"IPV6"`
	IPv4Mask          string         `db:"mask" csv:"IPV4_MASK"`
	IPv6Mask          string         `db:"-" csv:"IPV6_MASK"`
	BeginTime         string         `db:"contract_date" csv:"BEGIN_TIME"`
	EndTime           string         `db:"-" csv:"END_TIME"`
	LineObject        string         `db:"-" csv:"LINE_OBJECT"`
	LineCross         string         `db:"-" csv:"LINE_CROSS"`
	LineBlock         string         `db:"-" csv:"LINE_BLOCK"`
	LinePair          string         `db:"-" csv:"LINE_PAIR"`
	LineReserved      string         `db:"-" csv:"LINE_RESERVED"`
	LocType           string         `db:"-" csv:"LOC_TYPE"`
	LocLac            string         `db:"-" csv:"LOC_LAC"`
	LocCell           string         `db:"-" csv:"LOC_CELL"`
	LocTa             string         `db:"-" csv:"LOC_TA"`
	LocCellWireless   string         `db:"-" csv:"LOC_CELL_WIRELESS"`
	LocMac            string         `db:"-" csv:"LOC_MAC"`
	LocLatitude       string         `db:"-" csv:"LOC_LATITUDE"`
	LocLongitude      string         `db:"-" csv:"LOC_LONGITUDE"`
	LocProjectionType string         `db:"-" csv:"LOC_PROJECTION_TYPE"`
	RecordAction      string         `db:"-" csv:"RECORD_ACTION"`
	InternalID1       string         `db:"id" csv:"INTERNAL_ID1"`
	InternalID2       string         `db:"-" csv:"INTERNAL_ID2"`
}

func (a *AbonIdent) Render(db *sqlx.DB) (r []string, err error) { //
	var abons []AbonIdentRow //
	err = db.Select(&abons, `select u.uid, 
s.district_id,
pi.contract_date,
u.id,
pi.email,
INET_NTOA(dv.ip) as ip,
dh.mac
from 
users u 
JOIN dv_main dv ON dv.uid=u.uid
LEFT JOIN users_pi pi ON pi.uid=u.uid 
LEFT JOIN builds b ON b.id=pi.location_id
LEFT JOIN streets s ON s.id=b.street_id
LEFT JOIN bills bi ON u.bill_id=bi.id
LEFT JOIN companies c ON c.id=u.company_id
LEFT JOIN dhcphosts_hosts dh ON dh.uid=u.uid
JOIN tarif_plans tp ON tp.id=dv.tp_id
`)
	if err != nil {
		return nil, err
	}
	// for i := range abons {
	// 	abons[i].Calc()
	// }
	r = csv.MarshalCSV(abons, ";", "\"")
	return r, nil
}

func (a *AbonIdent) GetFileName() string {
	return fmt.Sprintf("ABONENT_IDENT_%s.txt", time.Now().Format("20060102_1504"))
}

func (a *AbonIdentRow) Calc() {
	a.IdentType = 5
	a.EquipmentType = 0
	a.IPType = 0
}
