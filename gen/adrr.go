package gen

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/csv"
)

type AbonAddrRow struct {
	AbonentID     string         `db:"-" csv:"ABONENT_ID"`
	RegionID      int            `db:"-" csv:"REGION_ID"`
	AddressTypeID int            `db:"-" csv:"ADDRESS_TYPE_ID"`
	AddressType   int            `db:"-" csv:"ADDRESS_TYPE"`
	Zip           sql.NullString `db:"zip" csv:"ZIP"`
	Country       string         `db:"-" csv:"COUNTRY"`
	Region        string         `db:"-" csv:"REGION"`
	Zone          string         `db:"-" csv:"ZONE"`
	City          sql.NullString `db:"dist" csv:"CITY"`
	Street        sql.NullString `db:"street" csv:"STREET"`
	Building      sql.NullString `db:"build" csv:"BUILDING"`
	BuildSect     string         `db:"-" csv:"BUILD_SECT"`
	Apartment     string         `db:"flat" csv:"APARTMENT"`
	UnstructInfo  string         `db:"-" csv:"UNSTRUCT_INFO"`
	BeginTime     time.Time      `db:"-" csv:"BEGIN_TIME" time:"2006-01-02 15:04:05"`
	EndTime       time.Time      `db:"-" csv:"END_TIME" time:"2006-01-02 15:04:05"`
	RecordAction  string         `db:"-" csv:"RECORD_ACTION"`
	InternalID1   string         `db:"uid" csv:"INTERNAL_ID1"`
	InternalID2   string         `db:"-" csv:"INTERNAL_ID2"`
	Descr         sql.NullString `db:"comments" csv:"-"`
	Company       int            `db:"company_id" csv:"-"`
	RegAddr       string         `db:"_regaddr" csv:"-"`
	RegAddrSame   bool           `db:"_regsame" csv:"-"`
}

type AbonAddr struct{}

func (a *AbonAddr) Render(db *sqlx.DB, cfg Config) (r []string, err error) { //
	var abons []AbonAddrRow //
	dta := cfg.InitDate.Format("2006-01-02")
	if cfg.OnlyOneDay {
		dta = time.Now().Format("2006-01-02")
	}
	err = db.Select(&abons, `SELECT u.uid, 
d.name as dist,
d.comments,
d.zip as zip,
s.name as street,
b.number as build,
pi.address_flat as flat,
u.company_id,
pi._regaddr,
pi._regsame
FROM 
users u 
JOIN dv_main dv ON dv.uid=u.uid
LEFT JOIN users_pi pi ON pi.uid=u.uid 
LEFT JOIN builds b ON b.id=pi.location_id
LEFT JOIN streets s ON s.id=b.street_id
LEFT JOIN districts d ON d.id=s.district_id
LEFT JOIN bills bi ON u.bill_id=bi.id
LEFT JOIN companies c ON c.id=u.company_id
JOIN tarif_plans tp ON tp.id=dv.tp_id
LEFT JOIN admin_actions aa1 on aa1.id = (select id from admin_actions 
	where uid=u.uid order by id limit 1)
WHERE aa1.datetime >= ?`, dta)
	if err != nil {
		return nil, err
	}
	var ainf addrInfo
	var rabons []AbonAddrRow
	for i := range abons {
		_ = json.Unmarshal([]byte(abons[i].Descr.String), &ainf)
		abons[i].Calc(ainf, cfg)
		rabons = append(rabons, abons[i])
		if abons[i].RegAddr != "" && !abons[i].RegAddrSame {
			ca := abons[i]
			ca.AddressTypeID = 0
			ca.AddressType = 1
			ca.Street = sql.NullString{}
			ca.Building = sql.NullString{}

			ca.UnstructInfo = abons[i].RegAddr
			rabons = append(rabons, ca)
		}
		if abons[i].RegAddrSame {
			ca := abons[i]
			ca.AddressTypeID = 0
			rabons = append(rabons, ca)
		}
	}
	r = csv.MarshalCSV(rabons, ";", "")
	return r, nil
}

func (a *AbonAddr) GetFileName() string {
	return fmt.Sprintf("ABONENT_ADDR_%s.txt", time.Now().Format("20060102_1504"))
}

type addrInfo struct {
	Country string `json:"country"`
	Region  string `json:"region"`
	Zone    string `json:"zone"`
}

func (a *AbonAddrRow) Calc(ainf addrInfo, cfg Config) {
	if a.Company > 0 {
		a.InternalID1 = fmt.Sprintf("%s%d", cfg.CompanyCode, a.Company)
	}
	a.Country = cfg.Country
	a.AddressTypeID = 3
	a.RegionID = cfg.RegionID
	a.BeginTime = cfg.InitDate
	a.Country = ainf.Country
	a.Zone = ainf.Zone
	a.Region = ainf.Region
	a.Building.String, a.BuildSect = SplitHouseNumber(a.Building.String)
	n, _ := strconv.Atoi(a.Building.String)
	if (!a.City.Valid && !a.Street.Valid) || n == 0 {
		a.AddressType = 1
		a.UnstructInfo = fmt.Sprintf("%s %s %s %s", a.City.String, a.Street.String, a.Building.String, a.Apartment)
		a.Building.String = ""
		a.BuildSect = ""
		a.Apartment = ""
	}
}

var LaterRegexp = regexp.MustCompile(`(\d+)[\ \/]?([а-яА-Я]$)`)
var CorpRegexp = regexp.MustCompile(`(\d+)(\D+)(\d+)`)

func SplitHouseNumber(house string) (build, sect string) {
	if house == "" {
		return "", ""
	}
	house = strings.TrimSpace(house)
	house = strings.Replace(house, `"`, "", -1)
	// for i := len(house) - 1; i >= 0; i-- {
	// 	if house[i] == '/' {
	// 		return house[:i], house[i+1:]
	// 	}
	// }
	if LaterRegexp.MatchString(house) {
		return LaterRegexp.FindStringSubmatch(house)[1],
			strings.ToUpper(LaterRegexp.FindStringSubmatch(house)[2])
	}
	if CorpRegexp.MatchString(house) {
		return CorpRegexp.FindStringSubmatch(house)[1],
			CorpRegexp.FindStringSubmatch(house)[2] + " " +
				strings.TrimSpace(CorpRegexp.FindStringSubmatch(house)[3])
	}

	return house, ""
}
