package gen

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/csv"
)

type Abons struct {
}

type AbonRow struct {
	AbonID               int            `db:"uid" csv:"ID"`                                    // ID
	RegionID             sql.NullInt64  `db:"district_id" csv:"REGION_ID"`                     // REGION_ID
	CDate                string         `db:"contract_date" csv:"CONTRACT_DATE"`               // CONTRACT_DATE
	Login                string         `db:"id" csv:"CONTRACT"`                               // CONTRACT
	BID                  int            `db:"bill_id" csv:"ACCOUNT"`                           // ACCOUNT
	ActualFrom           time.Time      `db:"-" csv:"ACTUAL_FROM" format:"2006-01-02"`         // ACTUAL_FROM
	ActualTo             time.Time      `db:"-" csv:"ACTUAL_TO" format:"2006-01-02"`           // ACTUAL_TO
	Company              int            `db:"company_id" csv:"-"`                              //
	AbonType             int            `db:"-" csv:"ABONENT_TYPE"`                            // ABONENT_TYPE
	NameInfoType         int            `db:"-" csv:"NAME_INFO_TYPE"`                          // NAME_INFO_TYPE
	FamilyName           string         `db:"-" csv:"FAMILY_NAME"`                             // FAMILY_NAME
	GivenName            string         `db:"-" csv:"GIVEN_NAME"`                              // GIVEN_NAME
	InitialName          string         `db:"-" csv:"INITIAL_NAME"`                            // INITIAL_NAME
	FIO                  string         `db:"fio" csv:"UNSTRUCT_NAME"`                         // UNSTRUCT_NAME
	BirthDate            time.Time      `db:"birth_date" csv:"BIRTH_DATE" format:"2006-01-02"` // BIRTH_DATE
	IdentCardTypeID      int            `db:"-" csv:"IDENT_CARD_TYPE_ID"`                      // IDENT_CARD_TYPE_ID
	IdentCardType        string         `db:"-" csv:"IDENT_CARD_TYPE"`                         // IDENT_CARD_TYPE
	IdentCardSerial      string         `db:"-" csv:"IDENT_CARD_SERIAL"`                       // IDENT_CARD_SERIAL
	IdentCardNumber      string         `db:"-" csv:"IDENT_CARD_NUMBER"`                       // IDENT_CARD_NUMBER
	IdentCardDescription string         `db:"-" csv:"IDENT_CARD_DESCRIPTION"`                  // IDENT_CARD_DESCRIPTION
	IdentCardUnstruct    string         `db:"passport" csv:"IDENT_CARD_UNSTRUCT"`              // IDENT_CARD_UNSTRUCT
	Bank                 string         `db:"-" csv:"BANK"`                                    // BANK
	BankAccount          string         `db:"-" csv:"BANK_ACCOUNT"`                            // BANK_ACCOUNT
	FullName             sql.NullString `db:"compname" csv:"FULL_NAME"`                        // FULL_NAME
	INN                  string         `db:"-" csv:"INN"`                                     // INN
	Contact              string         `db:"-" csv:"CONTACT"`                                 // CONTACT
	PhoneFax             string         `db:"phone" csv:"PHONE_FAX"`                           // PHONE_FAX
	Status               int            `db:"disdel" csv:"STATUS"`                             // STATUS
	Attach               string         `db:"-" csv:"ATTACH"`                                  // ATTACH
	Detach               string         `db:"-" csv:"DETACH"`                                  // DETACH
	NetworkType          int            `db:"-" csv:"NETWORK_TYPE"`                            // NETWORK_TYPE
	RecordAction         string         `db:"-" csv:"RECORD_ACTION"`                           // RECORD_ACTION
	InternalID1          string         `db:"uid" csv:"INTERNAL_ID1"`                          // INTERNAL_ID1
}

func (a *Abons) Render(db *sqlx.DB) (r []string, err error) { //
	var abons []AbonRow //
	err = db.Select(&abons, `select u.uid, 
s.district_id,
pi.contract_date,
u.id,
u.company_id,
pi.fio,
concat(pasport_num,' ',pasport_date,' ',pasport_grant) as passport,
c.name as compname,
u.disable+u.deleted as disdel
from 
users u 
JOIN dv_main dv ON dv.uid=u.uid
LEFT JOIN users_pi pi ON pi.uid=u.uid 
LEFT JOIN builds b ON b.id=pi.location_id
LEFT JOIN streets s ON s.id=b.street_id
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
	r = csv.MarshalCSV(abons, ";", "\"")
	return r, nil
}

func (a *Abons) GetFileName() string {
	return fmt.Sprintf("ABONENT_%s.txt", time.Now().Format("20060102_1504"))
}

func (r *AbonRow) Calc() {
	// r.ActualFrom = time.Now()
	// r.ActualTo = time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)
	r.AbonType = 42
	if r.Company > 0 {
		r.AbonType = 43
	}
	if r.Status == 0 {
		r.Status = 1
	} else {
		r.Status = 0
	}
	r.NameInfoType = 1
	r.IdentCardTypeID = 1
	r.IdentCardType = "Паспорт"
	r.InternalID1 = fmt.Sprintf("%d", r.AbonID)
	r.NetworkType = 4
}
