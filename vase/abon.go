package vase

import (
	"database/sql"
	"fmt"
	"net"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/config"
	"github.com/prgra/qgen/csv"
	"github.com/prgra/qgen/gen"
)

type AbonIdent struct{}

type AbonIdentRow struct {
	// Стандарт связи абонента по справочнику (GSM/CDMA/ТФОП/ШПД и т.д.).
	AbonType int `db:"-" csv:"abontype"`
	// Уникальный идентификатор абонента в выгрузке.
	ID string `db:"uid" csv:"id"`
	// Логин абонента для привязки к учетным данным/AAA.
	Login string `db:"login" csv:"login"`
	// Номер договора абонента.
	Dogovor string `db:"dogovor" csv:"dogovor"`
	// Текущий статус: 0 подключен, 1 отключен.
	Status int `db:"status" csv:"status"`
	// Дата и время заключения договора.
	ContractDate time.Time `db:"contract_date" csv:"contract_date" time:"02.01.2006T15:04:05"`
	// Дата и время расторжения договора.
	ContractEndDate time.Time `db:"contract_end_date" csv:"contract_end_date" time:"02.01.2006T15:04:05"`
	// Тип абонента: ФЛ, ЮЛ или оконечный пользователь ЮЛ.
	AbonentType int `db:"-" csv:"abonent_type"`
	// Тип представления ФИО: структурированное/неструктурированное.
	FIODataType int `db:"-" csv:"fio_data_type"`
	// Имя абонента (для структурированного ФИО).
	FirstName string `db:"-" csv:"first_name"`
	// Отчество абонента (для структурированного ФИО).
	MiddleName string `db:"-" csv:"middle_name"`
	// Фамилия абонента (для структурированного ФИО).
	LastName string `db:"-" csv:"last_name"`
	// Неструктурированное ФИО одной строкой.
	FIOUnstructured string `db:"-" csv:"fio_unstructured"`
	// Дата рождения абонента.
	BirthDate time.Time `db:"-" csv:"birth_date" time:"02.01.2006"`
	// Тип паспортных данных: структурированные/неструктурированные.
	PassportDataType int `db:"-" csv:"passport_data_type"`
	// Серия документа, удостоверяющего личность.
	DocSeries string `db:"-" csv:"doc_series"`
	// Номер документа, удостоверяющего личность.
	DocNumber         string       `db:"doc_number" csv:"doc_number"`
	PassportIssueDate sql.NullTime `db:"passport_issue_date" csv:"-"`
	// Дата выдачи документа.
	DocIssueDate time.Time `db:"-" csv:"doc_issue_date" time:"02.01.2006"`
	// Кем выдан документ.
	DocIssuedBy string `db:"doc_issued_by" csv:"doc_issued_by"`
	// Неструктурированные паспортные данные.
	PassportUnstruct string `db:"passport_unstruct" csv:"passport_unstruct"`
	// Тип документа (паспорт РФ/СССР и т.п.).
	DocType string `db:"-" csv:"doc_type"`
	// Банк абонента (для расчетов с оператором).
	BankName string `db:"-" csv:"bank_name"`
	// Номер банковского счета абонента.
	BankAccount string `db:"-" csv:"bank_account"`
	// Полное наименование юридического лица.
	ULName string `db:"-" csv:"ul_name"`
	// ИНН юридического лица.
	ULINN string `db:"-" csv:"ul_inn"`
	// Контактное лицо юридического лица.
	ULContact string `db:"-" csv:"ul_contact"`
	// Контактные телефоны юридического лица.
	ULContactPhones string `db:"-" csv:"ul_contact_phones"`
	// Банк юридического лица.
	ULBankName string `db:"-" csv:"ul_bank_name"`
	// Банковский счет юридического лица.
	ULBankAccount string `db:"-" csv:"ul_bank_account"`
	// Тип данных адреса регистрации: структурированный/неструктурированный.
	RegAddrDataType int `db:"-" csv:"reg_addr_data_type"`
	// Почтовый индекс адреса регистрации.
	RegAddrZIP string `db:"-" csv:"reg_addr_zip"`
	// Страна адреса регистрации.
	RegAddrCountry string `db:"-" csv:"reg_addr_country"`
	// Сокращенный тип региона адреса регистрации.
	RegAddrRegionType string `db:"-" csv:"reg_addr_region_type"`
	// Регион адреса регистрации.
	RegAddrRegion string `db:"-" csv:"reg_addr_region"`
	// Сокращенный тип района адреса регистрации.
	RegAddrDistrictType string `db:"-" csv:"reg_addr_district_type"`
	// Район/муниципальный округ адреса регистрации.
	RegAddrDistrict string `db:"-" csv:"reg_addr_district"`
	// Сокращенный тип населенного пункта адреса регистрации.
	RegAddrCityType string `db:"-" csv:"reg_addr_city_type"`
	// Населенный пункт адреса регистрации.
	RegAddrCity string `db:"-" csv:"reg_addr_city"`
	// Сокращенный тип улицы адреса регистрации.
	RegAddrStreetType string `db:"-" csv:"reg_addr_street_type"`
	// Улица адреса регистрации.
	RegAddrStreet string `db:"-" csv:"reg_addr_street"`
	// Номер дома адреса регистрации.
	RegAddrHouse string `db:"-" csv:"reg_addr_house"`
	// Корпус/строение адреса регистрации.
	RegAddrBuilding string `db:"-" csv:"reg_addr_building"`
	// Сокращенный тип помещения (кв/оф).
	RegAddrFlatType string `db:"-" csv:"reg_addr_flat_type"`
	// Номер квартиры/офиса адреса регистрации.
	RegAddrFlat string `db:"-" csv:"reg_addr_flat"`
	// Неструктурированный адрес регистрации.
	RegAddrUnstruct string `db:"-" csv:"reg_addr_unstruct"`
	// Тип данных адреса установки оборудования.
	InstAddrDataType int `db:"-" csv:"inst_addr_data_type"`
	// Почтовый индекс адреса установки.
	InstAddrZIP string `db:"-" csv:"inst_addr_zip"`
	// Страна адреса установки.
	InstAddrCountry string `db:"-" csv:"inst_addr_country"`
	// Сокращенный тип региона адреса установки.
	InstAddrRegionType string `db:"-" csv:"inst_addr_region_type"`
	// Регион адреса установки.
	InstAddrRegion string `db:"-" csv:"inst_addr_region"`
	// Сокращенный тип района адреса установки.
	InstAddrDistrictType string `db:"-" csv:"inst_addr_district_type"`
	// Район/муниципальный округ адреса установки.
	InstAddrDistrict string `db:"-" csv:"inst_addr_district"`
	// Сокращенный тип населенного пункта адреса установки.
	InstAddrCityType string `db:"-" csv:"inst_addr_city_type"`
	// Населенный пункт адреса установки.
	InstAddrCity string `db:"-" csv:"inst_addr_city"`
	// Сокращенный тип улицы адреса установки.
	InstAddrStreetType string `db:"-" csv:"inst_addr_street_type"`
	// Улица адреса установки.
	InstAddrStreet string `db:"-" csv:"inst_addr_street"`
	// Номер дома адреса установки.
	InstAddrHouse string `db:"-" csv:"inst_addr_house"`
	// Корпус/строение адреса установки.
	InstAddrBuilding string `db:"-" csv:"inst_addr_building"`
	// Сокращенный тип помещения адреса установки.
	InstAddrFlatType string `db:"-" csv:"inst_addr_flat_type"`
	// Номер квартиры/офиса адреса установки.
	InstAddrFlat string `db:"-" csv:"inst_addr_flat"`
	// Неструктурированный адрес установки.
	InstAddrUnstruct string `db:"-" csv:"inst_addr_unstruct"`
	// Тип данных почтового адреса абонента.
	PostAddrDataType int `db:"-" csv:"post_addr_data_type"`
	// Почтовый индекс почтового адреса.
	PostAddrZIP string `db:"-" csv:"post_addr_zip"`
	// Страна почтового адреса.
	PostAddrCountry string `db:"-" csv:"post_addr_country"`
	// Сокращенный тип региона почтового адреса.
	PostAddrRegionType string `db:"-" csv:"post_addr_region_type"`
	// Регион почтового адреса.
	PostAddrRegion string `db:"-" csv:"post_addr_region"`
	// Сокращенный тип района почтового адреса.
	PostAddrDistrictType string `db:"-" csv:"post_addr_district_type"`
	// Район/муниципальный округ почтового адреса.
	PostAddrDistrict string `db:"-" csv:"post_addr_district"`
	// Сокращенный тип населенного пункта почтового адреса.
	PostAddrCityType string `db:"-" csv:"post_addr_city_type"`
	// Населенный пункт почтового адреса.
	PostAddrCity string `db:"-" csv:"post_addr_city"`
	// Сокращенный тип улицы почтового адреса.
	PostAddrStreetType string `db:"-" csv:"post_addr_street_type"`
	// Улица почтового адреса.
	PostAddrStreet string `db:"-" csv:"post_addr_street"`
	// Номер дома почтового адреса.
	PostAddrHouse string `db:"-" csv:"post_addr_house"`
	// Корпус/строение почтового адреса.
	PostAddrBuilding string `db:"-" csv:"post_addr_building"`
	// Сокращенный тип помещения почтового адреса.
	PostAddrFlatType string `db:"-" csv:"post_addr_flat_type"`
	// Номер квартиры/офиса почтового адреса.
	PostAddrFlat string `db:"-" csv:"post_addr_flat"`
	// Неструктурированный почтовый адрес.
	PostAddrUnstruct string `db:"-" csv:"post_addr_unstruct"`
	// Тип данных адреса доставки счета.
	BillAddrDataType int `db:"-" csv:"bill_addr_data_type"`
	// Почтовый индекс адреса доставки счета.
	BillAddrZIP string `db:"-" csv:"bill_addr_zip"`
	// Страна адреса доставки счета.
	BillAddrCountry string `db:"-" csv:"bill_addr_country"`
	// Сокращенный тип региона адреса доставки счета.
	BillAddrRegionType string `db:"-" csv:"bill_addr_region_type"`
	// Регион адреса доставки счета.
	BillAddrRegion string `db:"-" csv:"bill_addr_region"`
	// Сокращенный тип района адреса доставки счета.
	BillAddrDistrictType string `db:"-" csv:"bill_addr_district_type"`
	// Район/муниципальный округ адреса доставки счета.
	BillAddrDistrict string `db:"-" csv:"bill_addr_district"`
	// Сокращенный тип населенного пункта адреса доставки счета.
	BillAddrCityType string `db:"-" csv:"bill_addr_city_type"`
	// Населенный пункт адреса доставки счета.
	BillAddrCity string `db:"-" csv:"bill_addr_city"`
	// Сокращенный тип улицы адреса доставки счета.
	BillAddrStreetType string `db:"-" csv:"bill_addr_street_type"`
	// Улица адреса доставки счета.
	BillAddrStreet string `db:"-" csv:"bill_addr_street"`
	// Номер дома адреса доставки счета.
	BillAddrHouse string `db:"-" csv:"bill_addr_house"`
	// Корпус/строение адреса доставки счета.
	BillAddrBuilding string `db:"-" csv:"bill_addr_building"`
	// Сокращенный тип помещения адреса доставки счета.
	BillAddrFlatType string `db:"-" csv:"bill_addr_flat_type"`
	// Номер квартиры/офиса адреса доставки счета.
	BillAddrFlat string `db:"-" csv:"bill_addr_flat"`
	// Неструктурированный адрес доставки счета.
	BillAddrUnstruct string `db:"-" csv:"bill_addr_unstruct"`
	// Начало интервала актуальности записи.
	ActualFrom time.Time `db:"actual_from" csv:"actual_from" time:"02.01.2006T15:04:05"`
	// Конец интервала актуальности записи.
	ActualTo time.Time `db:"actual_to" csv:"actual_to" time:"02.01.2006T15:04:05"`
	// SSID абонентской Wi-Fi точки доступа.
	WifiSSID string `db:"ssid" csv:"wifi_ssid"`
	// MAC-адрес абонентской Wi-Fi точки доступа.
	WifiMAC sql.NullString `db:"mac" csv:"wifi_mac"`
	// Контактные телефоны для взаимодействия с абонентом.
	ContactPhones string `db:"phone" csv:"contact_phones"`

	// Служебные поля выборки, в CSV не выгружаются.
	FIORaw     sql.NullString `db:"fio" csv:"-"`
	SBirthDate string         `db:"_birth_date" csv:"-"`
	City       sql.NullString `db:"city" csv:"-"`
	Dist       sql.NullString `db:"dist" csv:"-"`
	Country    sql.NullString `db:"country" csv:"-"`
	RegZIP     sql.NullString `db:"zip" csv:"-"`
	Street     sql.NullString `db:"street" csv:"-"`
	House      sql.NullString `db:"house" csv:"-"`
	Flat       sql.NullString `db:"flat" csv:"-"`
}

func (a *AbonIdent) Render(db *sqlx.DB, cfg config.Config) (r []string, err error) { //
	var abons []AbonIdentRow //
	dta := cfg.InitDate.Format("2006-01-02")
	if cfg.OnlyOneDay {
		dta = time.Now().Format("2006-01-02")
	}
	fmt.Println(dta)
	err = db.Select(&abons, `select
u.uid,
u.id as login,
pi.contract_id as dogovor,
u.deleted + u.disable as status,
pi.contract_date as contract_date,
aa2.datetime as contract_end_date,
aa1.datetime as actual_from,
aa2.datetime as actual_to,
dh.mac,
pi.phone,
pi.pasport_num as doc_number,
pi.pasport_date as passport_issue_date,
pi.pasport_grant as doc_issued_by,
pi.pasport_num as passport_unstruct,
pi.fio,
pi._birth_date,
pi.address_flat as flat,
d.zip as zip,
d.country as country,
d.name as dist,
d.city as city,
s.name as street,
b.number as house
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
LEFT JOIN districts d ON d.id=s.district_id
LEFT JOIN bills bi ON u.bill_id=bi.id
LEFT JOIN companies c ON c.id=u.company_id
LEFT JOIN dhcphosts_hosts dh ON dh.uid=u.uid
JOIN tarif_plans tp ON tp.id=dv.tp_id
WHERE aa2.datetime >= ?`, dta)
	if err != nil {
		return nil, err
	}
	for i := range abons {
		abons[i].Calc(cfg)
	}
	// плоский формат v3: все поля в кавычках, строка заголовка не выгружается
	r = csv.MarshalCSVNoHeader(abons, ";", "\"")
	return r, nil
}

func (a *AbonIdent) GetFileName() string {
	// subscribers_v3_DDmmYYYYHHMMSS.txt по документации VAS Experts
	return fmt.Sprintf("subscribers_v3_%s.txt", time.Now().Format("02012006150405"))
}

var passportNumberRe = regexp.MustCompile(`\D+`)
var phoneInputRe = regexp.MustCompile(`^\+?[0-9\s().-]+$`)

func (a *AbonIdentRow) Calc(cfg config.Config) {
	_ = cfg
	a.sanitize()
	if a.AbonType == 0 {
		a.AbonType = 4
	}
	if a.Status != 0 {
		a.Status = 1
	}
	a.ContactPhones = normalizePhone(a.ContactPhones)
	a.splitFIO()
	a.fillAddrs(cfg)
	a.fillBirthDate()
	if a.PassportIssueDate.Valid {
		a.DocIssueDate = a.PassportIssueDate.Time
	}
	passportNumber := string(passportNumberRe.ReplaceAll([]byte(a.DocNumber), nil))
	if len(passportNumber) == 10 && !a.DocIssueDate.IsZero() {
		a.PassportDataType = 0
		a.DocSeries = passportNumber[:4]
		a.DocNumber = passportNumber[4:]
		a.PassportUnstruct = ""
		a.DocType = "Паспорт РФ"
	} else if a.DocNumber != "" || a.DocIssuedBy != "" || !a.DocIssueDate.IsZero() || a.PassportUnstruct != "" {
		a.PassportDataType = 1
		a.DocSeries = ""
		a.DocNumber = ""
		a.DocIssueDate = time.Time{}
		a.DocIssuedBy = ""
	}
	a.WifiMAC.String = MakeMac(a.WifiMAC.String)
}

// sanitize заменяет переводы строк пробелами во всех строковых полях строки
// выгрузки: в плоском формате перенос внутри поля ломает построчный разбор.
func (a *AbonIdentRow) sanitize() {
	v := reflect.ValueOf(a).Elem()
	for i := 0; i < v.NumField(); i++ {
		switch f := v.Field(i); f.Kind() {
		case reflect.String:
			f.SetString(cleanNewlines(f.String()))
		case reflect.Struct:
			if ns, ok := f.Addr().Interface().(*sql.NullString); ok && ns.Valid {
				ns.String = cleanNewlines(ns.String)
			}
		}
	}
}

func cleanNewlines(s string) string {
	s = strings.ReplaceAll(s, "\r\n", " ")
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	return strings.TrimSpace(s)
}

func normalizePhone(phone string) string {
	if !phoneInputRe.MatchString(phone) {
		return ""
	}
	digits := string(passportNumberRe.ReplaceAll([]byte(phone), nil))
	if len(digits) == 11 && (digits[0] == '7' || digits[0] == '8') {
		digits = digits[1:]
	}
	if len(digits) != 10 {
		return ""
	}
	// по примеру из доки v3 телефоны выгружаются без плюса: 79211112233
	return "7" + digits
}

// splitFIO разбирает ФИО одной строкой "Фамилия Имя Отчество" на части.
func (a *AbonIdentRow) splitFIO() {
	parts := strings.Fields(a.FIORaw.String)
	switch {
	case len(parts) == 0:
		return
	case len(parts) == 1:
		a.FIODataType = 1
		a.FIOUnstructured = parts[0]
	default:
		a.FIODataType = 0
		a.LastName = parts[0]
		a.FirstName = parts[1]
		a.MiddleName = strings.Join(parts[2:], " ")
	}
}

// fillBirthDate парсит дату рождения из pi._birth_date (строка с произвольным форматом).
func (a *AbonIdentRow) fillBirthDate() {
	if a.SBirthDate == "" {
		return
	}
	t, err := gen.ParseDateFromString(a.SBirthDate)
	if err != nil {
		fmt.Println("parse birth date", err, a.ID, a.SBirthDate)
		return
	}
	a.BirthDate = t
}

// fillAddrs заполняет структурированные адрес регистрации и адрес установки.
func (a *AbonIdentRow) fillAddrs(cfg config.Config) {
	if a.City.String == "" && a.Dist.String == "" && a.Street.String == "" {
		return
	}
	country := a.Country.String
	if country == "" {
		country = cfg.Country
	}
	// город: d.city, при отсутствии - d.name (район)
	city := a.City.String
	district := a.Dist.String
	if city == "" && district != "" {
		city = district
		district = ""
	}
	house, corp := splitHouse(a.House.String)

	// адрес регистрации (по паспорту / месту подключения)
	a.RegAddrDataType = 0
	a.RegAddrZIP = a.RegZIP.String
	a.RegAddrCountry = country
	a.RegAddrDistrict = district
	a.RegAddrCity = city
	a.RegAddrStreet = a.Street.String
	a.RegAddrHouse = house
	a.RegAddrBuilding = corp
	a.RegAddrFlat = a.Flat.String
	// адрес установки оборудования совпадает с адресом подключения
	a.InstAddrDataType = a.RegAddrDataType
	a.InstAddrZIP = a.RegAddrZIP
	a.InstAddrCountry = a.RegAddrCountry
	a.InstAddrDistrict = a.RegAddrDistrict
	a.InstAddrCity = a.RegAddrCity
	a.InstAddrStreet = a.RegAddrStreet
	a.InstAddrHouse = a.RegAddrHouse
	a.InstAddrBuilding = a.RegAddrBuilding
	a.InstAddrFlat = a.RegAddrFlat
}

// splitHouse отделяет от номера дома корпус/литеру: "11А" -> "11", "А".
func splitHouse(house string) (build, corp string) {
	house = strings.TrimSpace(strings.ReplaceAll(house, `"`, ""))
	if house == "" {
		return "", ""
	}
	if m := houseCorpRe.FindStringSubmatch(house); m != nil {
		return m[1], strings.ToUpper(m[2])
	}
	if m := houseCorpNumRe.FindStringSubmatch(house); m != nil {
		return m[1], strings.ToUpper(strings.TrimSpace(m[3]))
	}
	return house, ""
}

var houseCorpRe = regexp.MustCompile(`(\d+)[ /-]?([а-яА-Я])$`)
var houseCorpNumRe = regexp.MustCompile(`(\d+)(\D+)(\d+)`)

// MakeMac - преобразует в строку вида 0A-0B-0C-0D-0E-0F
func MakeMac(mac string) string {
	pm, err := net.ParseMAC(mac)
	if err != nil {
		return ""
	}
	parts := make([]string, len(pm))
	for i, b := range pm {
		parts[i] = fmt.Sprintf("%02X", b)
	}
	return strings.Join(parts, "-")
}

// MakeIP - преобразует в строку вида 0A0B0C0D
func MakeIP(ip string) (hip string) {
	if ip == "" || ip == "0.0.0.0" {
		return ""
	}
	nip := net.ParseIP(ip)
	if nip == nil {
		return ""
	}
	return fmt.Sprintf("%0X", []byte(nip[12:]))
}

func (a *AbonIdent) GetRemoteDir() string {
	return ""
}
