package yhnt

import (
	"sort"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/config"
	"github.com/prgra/qgen/csv"
	"github.com/prgra/qgen/gen"
)

// PaymentRow is a row from the PAYMENTS table
type PaymentsTypeRow struct {
	DepID        int       // идентификатор филиала (число);
	PayType      int       // идентификатор типа оплаты (число);
	PayStartDate time.Time `time:"02.01.2006 15:04:05"` // дата начала действия типа оплаты (дата);
	PayEndDate   time.Time `time:"02.01.2006 15:04:05"` // дата завершения действия типа оплаты (в случае действия типа платежей - пустое значение);
	Descr        string    // описание/наименование типа оплаты (строка).

}

// Payments is a generator for PAYMENTS table
type PaymentsType struct{}

func (a *PaymentsType) Render(db *sqlx.DB, cfg config.Config) (r []string, err error) {
	_, pmap, err := gen.LoadPayMethodsMapFromFile("paymethods.map", cfg)
	if err != nil {
		return nil, err
	}
	var arr []int

	for k := range pmap {
		arr = append(arr, k)
	}
	sort.Ints(arr)
	var pays []PaymentsTypeRow
	for _, v := range pmap {

		pays = append(pays, PaymentsTypeRow{
			DepID:        cfg.RegionID,
			PayType:      v.ID,
			PayStartDate: cfg.InitDate,
			PayEndDate:   time.Time{},
			Descr:        v.Descr,
		})
	}

	r = csv.MarshalCSVNoHeader(pays, ";", "\"")
	return r, nil
}

func (a *PaymentsType) GetFileName() string {
	return "pay-types.csv"
}

func (a *PaymentsType) GetRemoteDir() string {
	return "dictionaries/pay-types"
}
