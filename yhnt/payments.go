package yhnt

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/config"
	"github.com/prgra/qgen/csv"
)

// PaymentRow is a row from the PAYMENTS table
type PaymentRow struct {
	DepID   int       // идентификатор филиала (число);
	Method  int       `db:"method"`                       // идентификатор типа оплаты (число);
	Login   string    `db:"login"`                        // 	логин или номер договора (строка);
	IP      string    `db:"ip"`                           // статический IP-адрес или IP-подсеть (при динамических адресах - не заполняется) (строка)
	DateU   int       `db:"date" csv:"-"`                 // unix для даты создания
	Date    time.Time `db:"-" time:"02.01.2006 15:04:05"` // дата и время пополнения баланса
	Sum     float32   `db:"sum"`                          // сумма перевода (строка).
	Comment string    `db:"dsc"`                          // сопутствующая информация о платеже, имеющаяся в распоряжении оператора связи (строка).
}

// Payments is a generator for PAYMENTS table
type Payments struct{}

func (a *Payments) Render(db *sqlx.DB, cfg config.Config) (r []string, err error) {
	var pays []PaymentRow
	dta := cfg.InitDate.Format("2006-01-02")
	if cfg.OnlyOneDay {
		dta = time.Now().Format("2006-01-02")
	}
	err = db.Select(&pays, `SELECT UNIX_TIMESTAMP(p.reg_date) as date, p.sum, u.id as login, INET_NTOA(dv.ip) as ip, p.dsc
		FROM payments p
		JOIN users u on u.uid = p.uid
		JOIN dv_main dv on dv.uid = p.uid
		where date >= ? order by p.id`, dta)
	if err != nil {
		return nil, err
	}
	for i := range pays {
		pays[i].Calc(cfg)
	}
	r = csv.MarshalCSVNoHeader(pays, ";", "\"")
	return r, nil
}

func (p *PaymentRow) Calc(cfg config.Config) {
	p.DepID = cfg.RegionID
	p.Date = time.Unix(int64(p.DateU), 0)

}

func (a *Payments) GetFileName() string {
	return "payments.csv"
}
