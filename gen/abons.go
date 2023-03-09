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

// ID
// REGION_ID
// CONTRACT_DATE
// CONTRACT
// ACCOUNT
// ACTUAL_FROM
// ACTUAL_TO
// ABONENT_TYPE
// NAME_INFO_TYPE
// FAMILY_NAME
// GIVEN_NAME
// INITIAL_NAME
// UNSTRUCT_NAME
// BIRTH_DATE
// IDENT_CARD_TYPE_ID
// IDENT_CARD_TYPE
// IDENT_CARD_SERIAL
// IDENT_CARD_NUMBER
// IDENT_CARD_DESCRIPTION
// IDENT_CARD_UNSTRUCT
// BANK
// BANK_ACCOUNT
// FULL_NAME
// INN
// CONTACT
// PHONE_FAX
// STATUS
// ATTACH
// DETACH
// NETWORK_TYPE
// RECORD_ACTION
// INTERNAL_ID1

type AbonRow struct {
	AbonID                 int           `db:"uid" csv:"ID"`
	RegID                  sql.NullInt64 `db:"district_id", csv:"REGION_ID"`
	CDate                  string        `db:"contract_date" csv:"CONTRACT_DATE"`
	Login                  string        `db:"id" csv:"CONTRACT"`
	BID                    int           `db:"bill_id" csv:"ACCOUNT"`
	Company                int           `db:"company_id"`
	AbonType               int           `db:"-" csv:"ABONENT_TYPE"`
	FIO                    string        `db:"fio"`
	BirthDate              time.Time     `db:"birth_date" csv:"BIRTH_DATE"`
	IDENT_CARD_TYPE_ID     int           `db:"-" csv:"IDENT_CARD_TYPE_ID"`
	IDENT_CARD_TYPE        string        `db:"-" csv:"IDENT_CARD_TYPE"`
	IDENT_CARD_SERIAL      string        `db:"-" csv:"IDENT_CARD_SERIAL"`
	IDENT_CARD_NUMBER      string        `db:"-" csv:"IDENT_CARD_NUMBER"`
	IDENT_CARD_DESCRIPTION string        `db:"-" csv:"IDENT_CARD_DESCRIPTION"`
	IDENT_CARD_UNSTRUCT    string        `db:"passport" csv:"IDENT_CARD_UNSTRUCT"`
	BANK                   string        `db:"-" csv:"BANK"`
	BANK_ACCOUNT           string        `db:"-" csv:"BANK_ACCOUNT"`
	FULL_NAME              string        `db:"-" csv:"FULL_NAME"`
	INN                    string        `db:"-" csv:"INN"`
	CONTACT                string        `db:"-" csv:"CONTACT"`
	PHONE_FAX              string        `db:"phone" csv:"PHONE_FAX"`
	STATUS                 string        `db:"-" csv:"STATUS"`
	ATTACH                 string        `db:"-" csv:"ATTACH"`
	DETACH                 string        `db:"-" csv:"DETACH"`
	NETWORK_TYPE           string        `db:"-" csv:"NETWORK_TYPE"`
	RECORD_ACTION          string        `db:"-" csv:"RECORD_ACTION"`
	INTERNAL_ID1           string        `db:"uid" csv:"INTERNAL_ID1"`
}

func (a *Abons) Render(db *sqlx.DB) (r []string, err error) {
	var abons []AbonRow
	err = db.Select(&abons, `select u.uid ,
s.district_id,
pi.contract_date,
u.id,
u.company_id,
pi.fio,
concat(pasport_num,pasport_date,pasport_grant) as passport
from 
users u 
JOIN dv_main dv ON dv.uid=u.uid
LEFT JOIN users_pi pi ON pi.uid=u.uid 
LEFT JOIN builds b ON b.id=pi.location_id
LEFT JOIN streets s ON s.id=b.street_id
LEFT JOIN bills bi ON u.bill_id=bi.id
JOIN tarif_plans tp ON tp.id=dv.tp_id
`)
	if err != nil {
		return nil, err
	}
	r = csv.MarshalCSV(abons, ";", "")
	return r, nil
}

func (a *Abons) GetFileName() string {
	return fmt.Sprintf("ABONENT_%s.txt", time.Now().Format("20060102_1504"))
}
