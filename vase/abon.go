package vase

import (
	"database/sql"
	"fmt"
	"net"
	"regexp"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/config"
	"github.com/prgra/qgen/csv"
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
COALESCE(pi.pasport, CONCAT_WS(', ', pi.pasport_num, pi.pasport_date, pi.pasport_grant)) as passport_unstruct
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
		return nil, err
	}
	for i := range abons {
		abons[i].Calc(cfg)
	}
	r = csv.MarshalCSV(abons, ";", "")
	return r, nil
}

func (a *AbonIdent) GetFileName() string {
	return fmt.Sprintf("ABONENT_IDENT_%s.txt", time.Now().Format("20060102_1504"))
}

var passportNumberRe = regexp.MustCompile(`\D+`)

func (a *AbonIdentRow) Calc(cfg config.Config) {
	_ = cfg
	if a.AbonType == 0 {
		a.AbonType = 4
	}
	if a.Status != 0 {
		a.Status = 1
	}
	if a.PassportIssueDate.Valid {
		a.DocIssueDate = a.PassportIssueDate.Time
	}
	passportNumber := string(passportNumberRe.ReplaceAll([]byte(a.DocNumber), nil))
	if len(passportNumber) == 10 && !a.DocIssueDate.IsZero() {
		a.PassportDataType = 0
		a.DocSeries = passportNumber[:4]
		a.DocNumber = passportNumber[4:]
		a.PassportUnstruct = ""
	} else if a.DocNumber != "" || a.DocIssuedBy != "" || !a.DocIssueDate.IsZero() || a.PassportUnstruct != "" {
		a.PassportDataType = 1
		a.DocSeries = ""
		a.DocNumber = ""
		a.DocIssueDate = time.Time{}
		a.DocIssuedBy = ""
	}
	a.WifiMAC.String = MakeMac(a.WifiMAC.String)
}

// MakeMac - преобразует в строку вида 0A0B0C0D0E0F
func MakeMac(mac string) string {
	pm, err := net.ParseMAC(mac)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%0X", []byte(pm))
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
