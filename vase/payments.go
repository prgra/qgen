package vase

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/config"
	"github.com/prgra/qgen/csv"
	"github.com/prgra/qgen/gen"
)

// Payments выгрузка платежей неопределённого вида payments_v2 (потоковая).
type Payments struct{}

type PaymentRow struct {
	// ID абонента, соответствует полю id выгрузки subscribers.
	ID string `db:"uid" csv:"id"`
	// Код абонента, зарезервировано на будущее - пустое поле.
	Code string `db:"-" csv:"code"`
	// Время и дата пополнения баланса.
	PaymentDate time.Time `db:"date" csv:"payment_date" time:"02.01.2006T15:04:05"`
	// Сумма платежа в рублях, формат 00000,00.
	Amount float64 `db:"sum" csv:"amount"`
	// Тип платежа: «Наличный» или «Безналичный».
	PayType string `db:"-" csv:"pay_type"`
	// Неструктурированная информация: платёжная система и идентификатор транзакции.
	PayParams string `db:"-" csv:"pay_params"`
	// ID физ лица, для плоского формата subscribers - пустое поле.
	PersonID string `db:"-" csv:"person_id"`

	// Служебные поля выборки, в CSV не выгружаются.
	PayID   int64  `db:"id" csv:"-"`
	Method  int    `db:"method" csv:"-"`
	Comment string `db:"dsc" csv:"-"`
}

func (a *Payments) Render(db *sqlx.DB, cfg config.Config) (r []string, err error) {
	var pays []PaymentRow
	dta := cfg.InitDate.Format("2006-01-02")
	if cfg.OnlyOneDay {
		dta = time.Now().Format("2006-01-02")
	}
	methods := loadPayMethods(cfg)
	err = db.Select(&pays, `SELECT p.id, p.uid, p.date, p.sum, p.method, COALESCE(p.dsc, '') as dsc
		FROM payments p
		JOIN users u on u.uid = p.uid
		where p.date >= ? and p.sum > 0 order by p.id`, dta)
	if err != nil {
		return nil, err
	}
	for i := range pays {
		pays[i].Calc(methods)
	}
	r = csv.MarshalCSVNoHeader(pays, ";", "\"")
	return r, nil
}

// loadPayMethods читает названия способов оплаты из paymethods.map (как в gen),
// без файла выгрузка работает, но без названия платёжной системы.
func loadPayMethods(cfg config.Config) map[int]string {
	r := make(map[int]string)
	_, pt, err := gen.LoadPayMethodsMapFromFile("paymethods.map", cfg)
	if err != nil {
		log.Println("paymethods.map", err)
		return r
	}
	for _, p := range pt {
		r[p.ID] = p.Descr
	}
	return r
}

func (p *PaymentRow) Calc(methods map[int]string) {
	// все даты выгружаются в UTC
	p.PaymentDate = p.PaymentDate.UTC()
	descr := strings.TrimSpace(methods[p.Method])
	p.PayType = "Безналичный"
	if p.Method == 0 || strings.Contains(strings.ToLower(descr), "налич") {
		p.PayType = "Наличный"
	}
	// для безналичного платежа обязательны платёжная система и идентификатор транзакции
	var params []string
	if descr != "" {
		params = append(params, descr)
	}
	params = append(params, fmt.Sprintf("ID платежа %d", p.PayID))
	if c := cleanNewlines(p.Comment); c != "" {
		params = append(params, c)
	}
	p.PayParams = strings.Join(params, ", ")
}

func (a *Payments) GetFileName() string {
	return fmt.Sprintf("payments_v2_%s.txt", time.Now().UTC().Format("02012006150405"))
}

func (a *Payments) GetRemoteDir() string {
	return ""
}
