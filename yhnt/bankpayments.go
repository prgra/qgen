package yhnt

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/config"
	"github.com/prgra/qgen/csv"
)

// BankPaymentRow is a row from the PAYMENTS table
type BankPaymentRow struct {
	DepID    int       // идентификатор филиала (число);
	Login    string    `db:"login"`                           // логин или номер договора (строка);
	AbonIP   string    `db:"ip"`                              // статический IP-адрес или IP-подсеть (при динамических адресах - не заполняется) (строка);
	Date     time.Time `db:"date" time:"02.01.2006 15:04:05"` // дата и время пополнения баланса (дата);
	BankFrom string    // номер банковского счета, с которого совершен платеж (строка);
	Bank     string    // наименование банка, со счета которого совершен перевод (строка);
	BankAddr string    // адрес банка, со счета которого совершен перевод (строка);
	Sum      string    `db:"sum"` // сумма перевода (строка).
}

// Payments is a generator for PAYMENTS table
type BankPayments struct{}

func (a *BankPayments) Render(db *sqlx.DB, cfg config.Config) (r []string, err error) {
	var pays []BankPaymentRow
	dta := cfg.InitDate.Format("2006-01-02")
	if cfg.OnlyOneDay {
		dta = time.Now().Format("2006-01-02")
	}
	err = db.Select(&pays, `SELECT p.date, p.sum, INET_NTOA(dv.ip) as ip, p.method, p.uid, p.bill_id, pi.phone 
		FROM payments p
		LEFT JOIN users_pi pi on pi.uid = p.uid
		LEFT JOIN dv_main dv on dv.uid = p.uid
		where date >= ? order by id`, dta)
	if err != nil {
		return nil, err
	}
	for i := range pays {
		pays[i].Calc(cfg)
	}
	r = csv.MarshalCSV(pays, ";", "")
	return r, nil
}

func (p *BankPaymentRow) Calc(cfg config.Config) {
	p.DepID = cfg.RegionID
}

func (a *BankPayments) GetFileName() string {
	return fmt.Sprintf("", time.Now().Format("20060102_1504"))
}

func (a *BankPayments) GetRemoteDir() string {
	return "payments/balance-fillup/"
}
