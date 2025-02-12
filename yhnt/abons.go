package yhnt

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/config"
	"github.com/prgra/qgen/csv"
	"github.com/prgra/qgen/gen"
)

type Abons struct{}

type AbonsRow struct {
	DepID       int            `db:"-"`                     //  1 идентификатор филиала (число)
	Login       string         `db:"login"`                 //  2 имя пользователя (логин для подключения к IP-сети) (строка, прочерк, если отсутствует)
	IP          string         `db:"ip"`                    //  3 статический IP-адрес или ip-подсеть (при динамических адресах - не заполняется) (строка)
	EMail       sql.NullString `db:"email"`                 //  4 адрес электронной почты (пустое поле «» если данных нет)
	Phone       string         `db:"phone"`                 //  5 номер телефона (пустое поле «» если данных нет)
	MacAddr     sql.NullString `db:"mac"`                   //  6 MAC-адрес абонента (при динамических адресах - не заполняется)
	CreateDateU int            `db:"crdate" csv:"-"`        // -- unix для даты создания
	CreateDate  time.Time      `time:"02.01.2006 15:04:05"` //  7 дата и время заключения договора (дата)
	ContractID  string         `db:"contract_id"`           //  8 номер договора (строка)
	Status      int            `db:"status"`                //  9 текущий статус абонента (0 - подключен, 1 - отключен) (число, «1» указывается при расторжении договора или, когда пользователь перестает пользоваться логином или статическим ip-адресом)
	EndDateU    sql.NullInt64  `db:"enddate" csv:"-"`       // -- unix для даты окончания
	ActualDate  time.Time      `time:"02.01.2006 15:04:05"` // 10 дата и время начала интервала времени, на котором актуальна информация (дата);
	EndDate     time.Time      `time:"02.01.2006 15:04:05"` // 11 дата и время окончания интервала времени, на котором актуальна информация (дата, обязательно заполняется при расторжении договора)
	Type        int            `db:"-"`                     // 12 тип абонента (0 - физическое лицо, 1 - юридическое лицо) (число)
	// информация об абоненте-физическом лице:
	FIOType int `db:"-"` // 13 тип данных по ФИО (0 - структурированные данные, 1 - неструктурированные) (число)
	// структурированное ФИО:
	SFIOName       string    `db:"-"`         // 14 имя (строка)
	SFIOPatronymic string    `db:"-"`         // 15 отчество (строка)
	SFIOSurname    string    `db:"-"`         // 16 фамилия (строка)
	USFIO          string    `db:"fio"`       // 17 неструктурированное ФИО (строка)
	BirthdayDate   time.Time `db:"_birthday"` // 18 дата рождения (дата)
	PassportType   int       `db:"-"`         // 19 тип паспортных данных (0 - структурированные паспортные данные, 1 - неструктурированные) (число)
	// структурированные паспортные данные:
	PasSeria   string    `db:""`                     // 19 серия удостоверения личности (строка);
	PasNumber  string    `db:"pasport_num"`          // 21 номер удостоверения личности (строка);
	SPasDate   string    `db:"pasport_date" csv:"-"` // -- для того чтобы распарсить дату ;
	PasDate    time.Time `db:"-" csv:"-"`            // -- дата ;
	PasVidano  string    `db:"pasport_grant"`        // 22 кем и когда выдано (строка);
	UnstuctPas string    `db:"pasport"`              // 23 неструктурированные паспортные данные (строка);
	DocType    int       `db:"-"`                    // 24 идентификатор типа документа, удостоверяющего личность (число);
	AbonBank   string    `db:"-"`                    // 25 банк абонента (используемый при расчете с оператором связи (строка), опциональное поле - заполняется в случае наличия таких сведений;
	BankAcc    string    `db:"-"`                    // 26 номер счета абонента в банке (используемый при расчетах с оператором связи) (строка), опциональное поле - заполняется в случае наличия таких сведений;
	// информация об абоненте-юридическом лице:
	CompanyID   int            `db:"company_id" csv:"-"` // -- идентификатор компании не выгружается в csv
	UrName      sql.NullString `db:"compname"`           // 27 полное наименование (строка);
	UrINN       sql.NullString `db:"tax_number"`         // 28 ИНН (строка);
	UrContact   sql.NullString `db:"representative"`     // 29 контактное лицо (строка);
	UrContPhone sql.NullString `db:"-"`                  // 30 контактные телефоны, факс (строка);
	UrBankName  sql.NullString `db:"bank_name"`          // 31 банк абонента, используемый при расчете с оператором связи (строка);
	UrBankSchet sql.NullString `db:"bank_account"`       // 32 номер счета абонента в банке, используемый при расчетах с оператором связи (строка);
	// адрес регистрации абонента (заполняется обязательно):
	AddrType int `db:"-"` // 33 тип данных по адресу регистрации абонента (0 - структурированные данные, 1 - неструктурированные) (число):
	// структурированный адрес:
	AddrZIP     sql.NullString `db:"zip"`     // 34 почтовый индекс, zip-код (строка);
	AddrCountry sql.NullString `db:"country"` // 25 страна (строка);
	AddrObl     string         `db:""`        // 36 область (строка);
	AddrDist    sql.NullString `db:"dist"`    // 37 район, муниципальный округ (строка);
	AddrCity    sql.NullString `db:"city"`    // 38 город/поселок/деревня/аул (строка);
	AddrStreet  sql.NullString `db:"street"`  // 39 улица (строка);
	AddrHouse   sql.NullString `db:"build"`   // 40 номер дома, строения (строка);
	AddrCorp    sql.NullString `db:"-"`       // 41 корпус (строка);
	AddFlat     sql.NullString `db:"flat"`    // 42 квартира, офис (строка);
	UnstAddr    sql.NullString `db:"addr"`    // 43 неструктурированный адрес (строка).
	// адрес устройства (заполняется обязательно):
	DevAddrType int `db:"-"` // 44 тип данных по адресу устройства (0 - структурированные данные, 1 - неструктурированные) (число)
	// структурированный адрес:
	DevAddrZIP     string `db:"-"` // 45 почтовый индекс, zip-код (строка);
	DevAddrCountry string `db:"-"` // 46 страна (строка);
	DevAddrObl     string `db:"-"` // 47 область (строка);
	DevAddrDist    string `db:"-"` // 48 район, муниципальный округ (строка);
	DevAddrCity    string `db:"-"` // 49 город/поселок/деревня/аул (строка);
	DevAddrStreet  string `db:"-"` // 50 улица (строка);
	DevAddrHouse   string `db:"-"` // 51 номер дома, строения (строка);
	DevAddrCorp    string `db:"-"` // 52 корпус (строка);
	DevAddFlat     string `db:"-"` // 53 квартира, офис (строка);
	DevUnstAddr    string `db:"-"` // 54 неструктурированный адрес (строка)
	// почтовый адрес (дополнительный адрес для юридических лиц):
	PostAddrType int `db:"-"` // 55 тип данных по почтовому адресу (0 - структурированные данные, 1 - неструктурированные) (число, пустое поле, если отсутствует):
	// структурированный адрес:
	PostAddrZIP     string // 56 почтовый индекс, zip-код (строка);
	PostAddrCountry string // 57 страна (строка);
	PostAddrObl     string // 58 область (строка);
	PostAddrDist    string // 59 район, муниципальный округ (строка);
	PostAddrCity    string // 60 город/поселок/деревня/аул (строка);
	PostAddrStreet  string // 61 улица (строка);
	PostAddrHouse   string // 62 номер дома, строения (строка);
	PostAddrCorp    string // 63 корпус (строка);
	PostAddFlat     string // 64 квартира, офис (строка);
	PostUnstAddr    string // 65 неструктурированный адрес (строка).
	// адрес доставки счета (дополнительный адрес для юридических лиц):
	DelivAddrType int `db:"-"` // 66 тип данных по адресу доставки счета (0 - структурированные данные, 1 - неструктурированные) (число, пустое поле, если отсутствует):
	// структурированный адрес:
	DelivAddrZIP     string // 67 почтовый индекс, zip-код (строка);
	DelivAddrCountry string // 68 страна (строка);
	DelivAddrObl     string // 69 область (строка);
	DelivAddrDist    string // 70 район, муниципальный округ (строка);
	DelivAddrCity    string // 71 город/поселок/деревня/аул (строка);
	DelivAddrStreet  string // 72 улица (строка);
	DelivAddrHouse   string // 73 номер дома, строения (строка);
	DelivAddrCorp    string // 74 корпус (строка);
	DelivAddFlat     string // 75 квартира, офис (строка);
	DelivUnstAddr    string // 76 неструктурированный адрес (строка).

}

var passportRe = regexp.MustCompile(`\D+`)

func (a *Abons) Render(db *sqlx.DB, cfg config.Config) (r []string, err error) { //
	var abons []AbonsRow //
	dta := cfg.InitDate.Format("2006-01-02")
	if cfg.OnlyOneDay {
		dta = time.Now().Format("2006-01-02")
	}
	err = db.Select(&abons, `select u.id as login,  
	pi.email,
u.company_id,
INET_NTOA(dv.ip) as ip,
dh.mac,
UNIX_TIMESTAMP(aa1.datetime) as crdate,
UNIX_TIMESTAMP(aa2.datetime) as enddate,
pi.fio,
pi.email,
pi.phone,
pi.contract_id,
u.deleted + u.disable as status,
c.name as compname,
c.tax_number,
c.bank_name,
c.representative,
c.bank_account,
d.zip,
d.country,
d.name as dist,
d.city,
s.name as street,
b.number as build,
pi.address_flat as flat,
pasport_num, pasport_date, pasport_grant,
concat(pasport_num,', ',pasport_date,', ',pasport_grant) as pasport
from 
users u 
JOIN dv_main dv ON dv.uid=u.uid
LEFT JOIN admin_actions aa1 on aa1.id = (select id from admin_actions 
	where uid=u.uid order by id limit 1)
LEFT JOIN admin_actions aa2 on aa2.id = (select id from admin_actions 
	where uid=u.uid order by id desc limit 1)

LEFT JOIN users_pi pi ON pi.uid=u.uid 
LEFT JOIN builds b ON b.id=pi.location_id
LEFT JOIN streets s ON s.id=b.street_id
LEFT JOIN districts d on d.id = s.district_id
LEFT JOIN bills bi ON u.bill_id=bi.id
LEFT JOIN companies c ON c.id=u.company_id
LEFT JOIN dhcphosts_hosts dh ON dh.uid=u.uid
JOIN tarif_plans tp ON tp.id=dv.tp_id
WHERE aa2.datetime >= ?`, dta)
	if err != nil {
		fmt.Println("err", err)
		return nil, err
	}
	for i := range abons {
		abons[i].Calc(cfg)
	}

	r = csv.MarshalCSVNoHeader(abons, ";", `"`)
	return r, nil
}

func (a *Abons) GetFileName() string {
	return fmt.Sprintf("abonents.csv")
}

func (a *AbonsRow) Calc(cfg config.Config) {
	a.DepID = cfg.RegionID
	a.CreateDate = time.Unix(int64(a.CreateDateU), 0)
	a.ActualDate = a.CreateDate
	a.FIOType = 1
	a.PassportType = 1
	a.MacAddr.String = strings.ToLower(a.MacAddr.String)
	a.EMail.String = strings.ToLower(a.EMail.String)
	if a.Status != 0 {
		a.Status = 1
		a.EndDate = time.Unix(int64(a.EndDateU.Int64), 0)
	}
	if a.IP == "0.0.0.0" {
		a.IP = "-"
	}
	if a.CompanyID > 0 || gen.IsUrLico(a.USFIO) {
		a.Type = 1
		a.UrName = sql.NullString{
			String: a.USFIO,
			Valid:  true,
		}
		a.USFIO = ""
	}
	n, k := gen.SplitHouseNumber(a.AddrHouse.String)
	a.AddrHouse = sql.NullString{
		String: n,
		Valid:  true,
	}
	a.AddrCorp = sql.NullString{
		String: k,
		Valid:  true,
	}
	a.PasDate, _ = gen.ParseDateFromString(a.SPasDate)
	pn := string(passportRe.ReplaceAll([]byte(a.PasNumber), []byte("")))
	if len(pn) == 10 && !a.PasDate.IsZero() {
		a.PasSeria = pn[:4]
		a.PasNumber = pn[4:]
		a.PasVidano = fmt.Sprintf("%s %s", a.PasVidano, a.PasDate.Format("02.01.2006"))
		a.UnstuctPas = ""
	} else {
		a.PasSeria = ""
		a.PasNumber = ""
		a.PasVidano = fmt.Sprintf("%s %s %s %s", a.PasSeria, a.PasNumber, a.PasVidano, a.PasDate.Format("02.01.2006"))
	}

	// a.DevAddFlat = a.AddFlat.String
	// a.DevAddrCity = a.AddrCity.String
	// a.DevAddrCountry = a.AddrCountry.String
	// a.DevAddrDist = a.AddrDist.String
	// a.DevAddrHouse = a.AddrHouse.String
	// a.DevAddrObl = a.AddrObl
	// a.DevAddrStreet = a.AddrStreet.String
	// a.DevAddrZIP = a.AddrZIP.String
	// a.DevUnstAddr = a.UnstAddr.String
	// a.DevAddrType = a.AddrType
	// a.PostAddrCity = a.AddrCity.String
	// a.PostAddrCountry = a.AddrCountry.String
	// a.PostAddrDist = a.AddrDist.String
	// a.PostAddrHouse = a.AddrHouse.String
	// a.PostAddrObl = a.AddrObl
	// a.PostAddrStreet = a.AddrStreet.String
	// a.PostAddrZIP = a.AddrZIP.String
	// a.PostUnstAddr = a.UnstAddr.String
	// a.PostAddrType = a.AddrType
	// a.DelivAddrCity = a.AddrCity.String
	// a.DelivAddrCountry = a.AddrCountry.String
	// a.DelivAddrDist = a.AddrDist.String
	// a.DelivAddrHouse = a.AddrHouse.String
	// a.DelivAddrObl = a.AddrObl
	// a.DelivAddrStreet = a.AddrStreet.String
	// a.DelivAddrZIP = a.AddrZIP.String
	// a.DelivUnstAddr = a.UnstAddr.String
	// a.DelivAddrType = a.AddrType
	// a.DelivAddrCorp = a.AddrCorp.String
	// a.DelivAddFlat = a.AddFlat.String
}

func (a *Abons) GetRemoteDir() string {
	return "abonents/abonents"
}
