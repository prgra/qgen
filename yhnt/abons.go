package yhnt

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/config"
	"github.com/prgra/qgen/csv"
	"github.com/prgra/qgen/gen"
)

type Abons struct{}

type AbonsRow struct {
	DepID       int            `db:"-"`                     // идентификатор филиала (число)
	Login       string         `db:"login"`                 // имя пользователя (логин для подключения к IP-сети) (строка, прочерк, если отсутствует)
	IP          string         `db:"ip"`                    // статический IP-адрес или ip-подсеть (при динамических адресах - не заполняется) (строка)
	EMail       sql.NullString `db:"email"`                 // адрес электронной почты (пустое поле «» если данных нет)
	Phone       string         `db:"phone"`                 // номер телефона (пустое поле «» если данных нет)
	MacAddr     sql.NullString `db:"mac"`                   // MAC-адрес абонента (при динамических адресах - не заполняется)
	CreateDateU int            `db:"crdate" csv:"-"`        // unix для даты создания
	CreateDate  time.Time      `time:"02.01.2006 15:04:05"` // дата и время заключения договора (дата)
	ContractID  string         `db:"contract_id"`           // номер договора (строка)
	Status      int            `db:"status"`                // текущий статус абонента (0 - подключен, 1 - отключен) (число, «1» указывается при расторжении договора или, когда пользователь перестает пользоваться логином или статическим ip-адресом)
	EndDateU    sql.NullInt64  `db:"enddate" csv:"-"`       // unix для даты окончания
	EndDate     time.Time      `time:"02.01.2006 15:04:05"` // дата и время окончания интервала времени, на котором актуальна информация (дата, обязательно заполняется при расторжении договора)
	Type        int            `db:"-"`                     // тип абонента (0 - физическое лицо, 1 - юридическое лицо) (число)
	// информация об абоненте-физическом лице:
	FIOType int `db:"-"` // тип данных по ФИО (0 - структурированные данные, 1 - неструктурированные) (число)
	// структурированное ФИО:
	SFIOName       string    `db:"-"`         // имя (строка)
	SFIOPatronymic string    `db:"-"`         // отчество (строка)
	SFIOSurname    string    `db:"-"`         // фамилия (строка)
	USFIO          string    `db:"fio"`       // неструктурированное ФИО (строка)
	BirthdayDate   time.Time `db:"_birthday"` // дата рождения (дата)
	PassportType   int       `db:"-"`         // тип паспортных данных (0 - структурированные паспортные данные, 1 - неструктурированные) (число)
	// структурированные паспортные данные:
	SPassSeria  string `db:"-"`     // серия удостоверения личности (строка);
	SPassNumber string `db:"-"`     //  номер удостоверения личности (строка);
	SpassVidano string `db:"-"`     // кем и когда выдано (строка);
	UnstuctPass string `db:"passp"` //  неструктурированные паспортные данные (строка);
	DocType     int    `db:"-"`     // идентификатор типа документа, удостоверяющего личность (число);
	AbonBank    string `db:"-"`     // банк абонента (используемый при расчете с оператором связи (строка), опциональное поле - заполняется в случае наличия таких сведений;
	BankAcc     string `db:"-"`     // номер счета абонента в банке (используемый при расчетах с оператором связи) (строка), опциональное поле - заполняется в случае наличия таких сведений;
	// информация об абоненте-юридическом лице:
	CompanyID   int    `db:"company_id" csv:"-"` // идентификатор компании не выгружается в csv
	UrName      string `db:"urname"`             // полное наименование (строка);
	UrINN       string `db:"tax_number"`         // ИНН (строка);
	UrContact   string `db:"representative"`     //  контактное лицо (строка);
	UrContPhone string `db:"-"`                  // контактные телефоны, факс (строка);
	UrBankName  string `db:"bank_name"`          // банк абонента, используемый при расчете с оператором связи (строка);
	UrBankSchet string `db:"bank_account"`       // номер счета абонента в банке, используемый при расчетах с оператором связи (строка);
	// адрес регистрации абонента (заполняется обязательно):
	AddrType int `db:"-"` // тип данных по адресу регистрации абонента (0 - структурированные данные, 1 - неструктурированные) (число):
	// структурированный адрес:
	AddrZIP     string `db:"-"`    // почтовый индекс, zip-код (строка);
	AddrCountry string `db:"-"`    // страна (строка);
	AddrObl     string `db:"-"`    // область (строка);
	AddrDist    string `db:"-"`    // район, муниципальный округ (строка);
	AddrCity    string `db:"-"`    // город/поселок/деревня/аул (строка);
	AddrStreet  string `db:"-"`    // улица (строка);
	AddrHouse   string `db:"-"`    // номер дома, строения (строка);
	AddrCorp    string `db:"-"`    // корпус (строка);
	AddFlat     string `db:"-"`    // квартира, офис (строка);
	UnstAddr    string `db:"addr"` // неструктурированный адрес (строка).
	// адрес устройства (заполняется обязательно):
	DevAddrType int `db:"-"` // тип данных по адресу устройства (0 - структурированные данные, 1 - неструктурированные) (число)
	// структурированный адрес:
	DevAddrZIP     string `db:"-"` // почтовый индекс, zip-код (строка);
	DevAddrCountry string `db:"-"` // страна (строка);
	DevAddrObl     string `db:"-"` // область (строка);
	DevAddrDist    string `db:"-"` // район, муниципальный округ (строка);
	DevAddrCity    string `db:"-"` // город/поселок/деревня/аул (строка);
	DevAddrStreet  string `db:"-"` // улица (строка);
	DevAddrHouse   string `db:"-"` // номер дома, строения (строка);
	DevAddrCorp    string `db:"-"` // корпус (строка);
	DevAddFlat     string `db:"-"` // квартира, офис (строка);
	DevUnstAddr    string `db:"-"` // неструктурированный адрес (строка)
	// почтовый адрес (дополнительный адрес для юридических лиц):
	PostAddrType int `db:"-"` // тип данных по почтовому адресу (0 - структурированные данные, 1 - неструктурированные) (число, пустое поле, если отсутствует):
	// структурированный адрес:
	PostAddrZIP     string // почтовый индекс, zip-код (строка);
	PostAddrCountry string // страна (строка);
	PostAddrObl     string // область (строка);
	PostAddrDist    string // район, муниципальный округ (строка);
	PostAddrCity    string // город/поселок/деревня/аул (строка);
	PostAddrStreet  string // улица (строка);
	PostAddrHouse   string // номер дома, строения (строка);
	PostAddrCorp    string // корпус (строка);
	PostAddFlat     string // квартира, офис (строка);
	PostUnstAddr    string // неструктурированный адрес (строка).
	// адрес доставки счета (дополнительный адрес для юридических лиц):
	DelivAddrType int `db:"-"` // тип данных по адресу доставки счета (0 - структурированные данные, 1 - неструктурированные) (число, пустое поле, если отсутствует):
	// структурированный адрес:
	DelivAddrZIP     string // почтовый индекс, zip-код (строка);
	DelivAddrCountry string // страна (строка);
	DelivAddrObl     string // область (строка);
	DelivAddrDist    string // район, муниципальный округ (строка);
	DelivAddrCity    string // город/поселок/деревня/аул (строка);
	DelivAddrStreet  string // улица (строка);
	DelivAddrHouse   string // номер дома, строения (строка);
	DelivAddrCorp    string // корпус (строка);
	DelivAddFlat     string // квартира, офис (строка);
	DelivUnstAddr    string // неструктурированный адрес (строка).

}

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
c.bank_account
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

	r = csv.MarshalCSV(abons, ";", `"`)
	return r, nil
}

func (a *Abons) GetFileName() string {
	return fmt.Sprintf("abonents_new.csv")
}

func (a *AbonsRow) Calc(cfg config.Config) {
	a.DepID = cfg.RegionID
	a.CreateDate = time.Unix(int64(a.CreateDateU), 0)
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
		a.UrName = a.USFIO
		a.USFIO = ""
	}
}
