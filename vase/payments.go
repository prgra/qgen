package vase

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/config"
	"github.com/prgra/qgen/csv"
)

// Payments выгрузка платежей неопределённого вида (ValidateBalanceFillupRecord).
type Payments struct{}

type PaymentRow struct {
	// Уникальный идентификатор абонента, как в поле id выгрузки subscribers.
	ID string `db:"uid" csv:"id"`
	// Логин абонента, как в поле login выгрузки subscribers.
	Login string `db:"login" csv:"login"`
	// Дата и время платежа (UTC).
	PaymentDate time.Time `db:"date" csv:"payment_date" time:"02.01.2006T15:04:05"`
	// Сумма платежа.
	Amount float64 `db:"sum" csv:"amount"`
	// Сопутствующая информация о платеже.
	Comment string `db:"dsc" csv:"comment"`
}

func (a *Payments) Render(db *sqlx.DB, cfg config.Config) (r []string, err error) {
	var pays []PaymentRow
	dta := cfg.InitDate.Format("2006-01-02")
	if cfg.OnlyOneDay {
		dta = time.Now().Format("2006-01-02")
	}
	err = db.Select(&pays, `SELECT p.uid, u.id as login, p.date, p.sum, COALESCE(p.dsc, '') as dsc
		FROM payments p
		JOIN users u on u.uid = p.uid
		where p.date >= ? order by p.id`, dta)
	if err != nil {
		return nil, err
	}
	for i := range pays {
		pays[i].Calc()
	}
	// формат как у subscribers v3: все поля в кавычках, строка заголовка не выгружается
	r = csv.MarshalCSVNoHeader(pays, ";", "\"")
	return r, nil
}

func (p *PaymentRow) Calc() {
	// все даты выгружаются в UTC
	p.PaymentDate = p.PaymentDate.UTC()
	p.Comment = cleanNewlines(p.Comment)
}

func (a *Payments) GetFileName() string {
	return fmt.Sprintf("payments_%s.txt", time.Now().UTC().Format("02012006150405"))
}

func (a *Payments) GetRemoteDir() string {
	return ""
}
