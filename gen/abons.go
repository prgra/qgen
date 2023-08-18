package gen

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/csv"
)

type Abons struct {
}

type AbonRow struct {
	AbonID               string         `db:"-" csv:"ID"`                                                   // ID
	RegionID             int            `db:"-" csv:"REGION_ID"`                                            // REGION_ID
	CDate                time.Time      `db:"contract_date" csv:"CONTRACT_DATE" time:"2006-01-02 15:04:05"` // CONTRACT_DATE
	Login                string         `db:"id" csv:"CONTRACT"`                                            // CONTRACT
	BID                  int            `db:"bill_id" csv:"ACCOUNT"`                                        // ACCOUNT
	ActualFrom           time.Time      `db:"-" csv:"ACTUAL_FROM" time:"2006-01-02 15:04:05"`               // ACTUAL_FROM
	ActualTo             time.Time      `db:"-" csv:"ACTUAL_TO" time:"2006-01-02 15:04:05"`                 // ACTUAL_TO
	Company              int            `db:"company_id" csv:"-"`                                           //
	AbonType             int            `db:"-" csv:"ABONENT_TYPE"`                                         // ABONENT_TYPE
	NameInfoType         int            `db:"-" csv:"NAME_INFO_TYPE"`                                       // NAME_INFO_TYPE
	FamilyName           string         `db:"-" csv:"FAMILY_NAME"`                                          // FAMILY_NAME
	GivenName            string         `db:"-" csv:"GIVEN_NAME"`                                           // GIVEN_NAME
	InitialName          string         `db:"-" csv:"INITIAL_NAME"`                                         // INITIAL_NAME
	FIO                  string         `db:"fio" csv:"UNSTRUCT_NAME"`                                      // UNSTRUCT_NAME
	BirthDate            time.Time      `db:"-" csv:"BIRTH_DATE" time:"2006-01-02"`                         // BIRTH_DATE
	SBirthDate           string         `db:"_birth_date" csv:"-"`                                          // BIRTH_DATE
	IdentCardTypeID      int            `db:"-" csv:"IDENT_CARD_TYPE_ID"`                                   // IDENT_CARD_TYPE_ID
	IdentCardType        int            `db:"-" csv:"IDENT_CARD_TYPE"`                                      // IDENT_CARD_TYPE
	IdentCardSerial      string         `db:"-" csv:"IDENT_CARD_SERIAL"`                                    // IDENT_CARD_SERIAL
	SPassportDate        time.Time      `db:"pasport_date" csv:"-"`                                         // SPassportDate для формирования паспорта
	IdentCardNumber      string         `db:"pasport_num" csv:"IDENT_CARD_NUMBER"`                          // IDENT_CARD_NUMBER
	IdentCardDescription string         `db:"pasport_grant" csv:"IDENT_CARD_DESCRIPTION"`                   // IDENT_CARD_DESCRIPTION
	IdentCardUnstruct    string         `db:"pasport" csv:"IDENT_CARD_UNSTRUCT"`                            // IDENT_CARD_UNSTRUCT
	Bank                 string         `db:"-" csv:"BANK"`                                                 // BANK
	BankAccount          string         `db:"-" csv:"BANK_ACCOUNT"`                                         // BANK_ACCOUNT
	FullName             sql.NullString `db:"-" csv:"FULL_NAME"`                                            // FULL_NAME
	INN                  string         `db:"-" csv:"INN"`                                                  // INN
	Contact              string         `db:"-" csv:"CONTACT"`                                              // CONTACT
	PhoneFax             string         `db:"phone" csv:"PHONE_FAX"`                                        // PHONE_FAX
	Status               int            `db:"disdel" csv:"STATUS"`                                          // STATUS
	Attach               sql.NullTime   `db:"attach" csv:"ATTACH" time:"2006-01-02 15:04:05"`               // ATTACH
	Detach               sql.NullTime   `db:"detach" csv:"DETACH" time:"2006-01-02 15:04:05"`               // DETACH
	NetworkType          int            `db:"-" csv:"NETWORK_TYPE"`                                         // NETWORK_TYPE
	RecordAction         string         `db:"-" csv:"RECORD_ACTION"`                                        // RECORD_ACTION
	InternalID1          string         `db:"uid" csv:"INTERNAL_ID1"`                                       // INTERNAL_ID1
	CompanyName          sql.NullString `db:"compname" csv:"-"`
}

func (a *Abons) Render(db *sqlx.DB) (r []string, err error) { //
	var abons []AbonRow //
	dta := EnvInitDate.Format("2006-01-02")
	if EnvOnlyOneDay {
		dta = time.Now().Format("2006-01-02")
	}
	err = db.Select(&abons, `select u.uid, 
pi.contract_date,
u.id,
u.company_id,
pi.fio,
u.bill_id,
pasport_num, pasport_date, pasport_grant,
concat(pasport_num,', ',pasport_date,', ',pasport_grant) as pasport,
c.name as compname,
u.disable+u.deleted as disdel,
aa1.datetime as attach,
aa2.datetime as detach
from 
users u 
LEFT JOIN admin_actions aa1 on aa1.id = (select id from admin_actions 
	where uid=u.uid order by id limit 1)
LEFT JOIN admin_actions aa2 on aa2.id = (select id from admin_actions 
	where uid=u.uid order by id desc limit 1)
JOIN dv_main dv ON dv.uid=u.uid
LEFT JOIN users_pi pi ON pi.uid=u.uid 
LEFT JOIN builds b ON b.id=pi.location_id
LEFT JOIN streets s ON s.id=b.street_id
LEFT JOIN bills bi ON u.bill_id=bi.id
LEFT JOIN companies c ON c.id=u.company_id
JOIN tarif_plans tp ON tp.id=dv.tp_id
WHERE aa1.datetime >= ?`, dta)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	additional := make(map[int]AbonRow)
	for i := range abons {
		abons[i].UrFizCalc()
		if abons[i].Company > 0 {
			r := abons[i]
			r.AbonType = 42
			key := fmt.Sprintf("%s%d", EnvCompanyCode, r.Company)
			r.InternalID1 = key
			abons[i].InternalID1 = key
			r.Company = 0
			r.Calc()
			additional[r.Company] = r
		}
		abons[i].Calc()

	}
	for k := range additional {
		abons = append(abons, additional[k])
	}
	r = csv.MarshalCSV(abons, ";", "")
	return r, nil
}

func (a *Abons) GetFileName() string {
	return fmt.Sprintf("ABONENT_%s.txt", time.Now().Format("20060102_1504"))
}

var passportRe = regexp.MustCompile(`\D+`)

func (r *AbonRow) UrFizCalc() {
	r.AbonType = 42
	if r.Company > 0 || IsUrLico(r.FIO) {
		r.AbonType = 43
	}
}

func (r *AbonRow) Calc() {
	r.ActualFrom = r.Attach.Time
	if r.AbonType == 43 {
		if r.CompanyName.Valid {
			r.FullName = r.CompanyName
		} else {
			r.FullName.String = r.FIO
			r.FullName.Valid = true
		}
		r.FIO = ""
	}
	if r.Status == 0 {
		r.Status = 0
		r.Detach = sql.NullTime{}
	} else {
		r.Status = 1
	}
	r.NameInfoType = 1
	r.IdentCardTypeID = 1
	r.IdentCardType = 1
	// r.InternalID1 = fmt.Sprintf("%d", r.AbonID)
	r.NetworkType = 4
	if r.AbonType == 42 {
		pn := string(passportRe.ReplaceAll([]byte(r.IdentCardNumber), []byte("")))
		if len(pn) == 10 && !r.SPassportDate.IsZero() {
			r.IdentCardType = 0
			r.IdentCardSerial = pn[:4]
			r.IdentCardNumber = pn[4:]
			r.IdentCardDescription = fmt.Sprintf("%s %s", r.SPassportDate.Format("2006-01-02"), r.IdentCardDescription)
			r.IdentCardUnstruct = ""
		} else {
			r.IdentCardSerial = ""
			r.IdentCardNumber = ""
			r.IdentCardDescription = ""
			r.IdentCardTypeID = 2
		}
	} else {
		r.IdentCardSerial = ""
		r.IdentCardNumber = ""
		r.IdentCardDescription = ""
		r.IdentCardTypeID = 0
		r.IdentCardUnstruct = ""
	}
	r.RegionID = EnvRegionID
	if r.CDate.IsZero() {
		r.CDate = EnvInitDate
	}
	// if r.Company > 0 {
	// 	r.InternalID1 = strconv.Itoa(r.Company)
	// 	r.AbonID = 0
	// }
}

func IsUrLico(s string) bool {
	if s == "" {
		return false
	}
	if strings.HasPrefix(s, "ОАО") ||
		strings.HasPrefix(s, "ЗАО") ||
		strings.HasPrefix(s, "ООО") ||
		strings.HasPrefix(s, "ПАО") ||
		strings.HasPrefix(s, "НКО") ||
		strings.HasPrefix(s, "НП") ||
		strings.HasPrefix(s, "АО") ||
		strings.HasPrefix(s, "АНО") ||
		strings.HasPrefix(s, "АКБ") ||
		strings.HasPrefix(s, "АК") {
		return true
	}
	if strings.Contains(s, `"`) {
		return true
	}
	return false
}
