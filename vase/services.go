package vase

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/config"
	"github.com/prgra/qgen/csv"
)

// Услуга доступа в интернет - единственная в номенклатуре (тарифы не выгружаются).
const (
	inetServiceID       = 1
	inetServiceMnemonic = "Интернет"
	inetServiceDescr    = "Высокоскоростной доступ в интернет"
)

// ServiceList выгрузка номенклатуры услуг оператора service_list_v1 (объектная).
type ServiceList struct{}

type ServiceListRow struct {
	// ID услуги, целое число от 1 до 4294967295.
	ID int `csv:"id"`
	// Мнемоническое обозначение (название услуги).
	Mnemonic string `csv:"mnemonic"`
	// Описание услуги.
	Description string `csv:"description"`
	// Дата активации услуги.
	BeginTime time.Time `csv:"begin_time" time:"02.01.2006"`
	// Дата деактивации услуги.
	EndTime time.Time `csv:"end_time" time:"02.01.2006"`
	// Тип сети: «Ethernet» - стационарная сеть передачи данных.
	NetworkType string `csv:"network_type"`
}

func (a *ServiceList) Render(db *sqlx.DB, cfg config.Config) (r []string, err error) {
	_ = db
	list := []ServiceListRow{
		{
			ID:          inetServiceID,
			Mnemonic:    inetServiceMnemonic,
			Description: inetServiceDescr,
			BeginTime:   cfg.InitDate,
			NetworkType: "Ethernet",
		},
	}
	r = csv.MarshalCSVNoHeader(list, ";", "\"")
	return r, nil
}

func (a *ServiceList) GetFileName() string {
	return fmt.Sprintf("service_list_v1_%s.txt", time.Now().UTC().Format("02012006150405"))
}

func (a *ServiceList) GetRemoteDir() string {
	return ""
}

// Services выгрузка услуг, выданных абонентам, services_v1 (объектная).
type Services struct{}

type ServiceRow struct {
	// ID услуги из номенклатуры service_list.
	ServiceID int `db:"-" csv:"service_id"`
	// ID абонента, соответствует полю id выгрузки subscribers.
	ID string `db:"uid" csv:"id"`
	// ID лицевого счёта, зарезервировано на будущее - пустое поле.
	AccountID string `db:"-" csv:"account_id"`
	// ID точки подключения, зарезервировано на будущее - пустое поле.
	PointID string `db:"-" csv:"point_id"`
	// ID договора, зарезервировано на будущее - пустое поле.
	ContractID string `db:"-" csv:"contract_id"`
	// Индивидуальные параметры настройки услуги абонента.
	Parameter string `db:"-" csv:"parameter"`
	// Дата подключения услуги, равна дате заключения договора.
	BeginTime time.Time `db:"-" csv:"begin_time" time:"02.01.2006T15:04:05"`
	// Дата отключения услуги, больше даты подключения.
	EndTime time.Time `db:"-" csv:"end_time" time:"02.01.2006T15:04:05"`
	// ID физ лица, для плоского формата subscribers - пустое поле.
	PersonID string `db:"-" csv:"person_id"`

	// Служебные поля выборки, в CSV не выгружаются.
	Status       int          `db:"status" csv:"-"`
	ContractDate sql.NullTime `db:"contract_date" csv:"-"`
	FirstAction  sql.NullTime `db:"actual_from" csv:"-"`
	LastAction   sql.NullTime `db:"actual_to" csv:"-"`
}

func (a *Services) Render(db *sqlx.DB, cfg config.Config) (r []string, err error) {
	var rows []ServiceRow
	dta := cfg.InitDate.Format("2006-01-02")
	if cfg.OnlyOneDay {
		dta = time.Now().Format("2006-01-02")
	}
	// тот же набор абонентов, что и в subscribers
	err = db.Select(&rows, `select
u.uid,
u.deleted + u.disable as status,
pi.contract_date as contract_date,
aa1.datetime as actual_from,
aa2.datetime as actual_to
from
users u
JOIN dv_main dv ON dv.uid=u.uid
LEFT JOIN admin_actions aa1 on aa1.id = (select id from admin_actions
	where uid=u.uid order by id limit 1)
LEFT JOIN admin_actions aa2 on aa2.id = (select id from admin_actions
	where uid=u.uid order by id desc limit 1)
LEFT JOIN users_pi pi ON pi.uid=u.uid
JOIN tarif_plans tp ON tp.id=dv.tp_id
WHERE aa2.datetime >= ?`, dta)
	if err != nil {
		return nil, err
	}
	for i := range rows {
		rows[i].Calc()
	}
	r = csv.MarshalCSVNoHeader(rows, ";", "\"")
	return r, nil
}

func (s *ServiceRow) Calc() {
	s.ServiceID = inetServiceID
	s.Parameter = inetServiceDescr
	switch {
	case s.ContractDate.Valid && !s.ContractDate.Time.IsZero():
		s.BeginTime = s.ContractDate.Time
	case s.FirstAction.Valid:
		s.BeginTime = s.FirstAction.Time
	}
	// дата отключения только у отключенных абонентов и строго позже подключения
	if s.Status != 0 && s.LastAction.Valid && s.LastAction.Time.After(s.BeginTime) {
		s.EndTime = s.LastAction.Time
	}
}

func (a *Services) GetFileName() string {
	return fmt.Sprintf("services_v1_%s.txt", time.Now().UTC().Format("02012006150405"))
}

func (a *Services) GetRemoteDir() string {
	return ""
}
