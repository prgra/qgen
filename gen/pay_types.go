package gen

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/csv"
)

type PayTypes struct{}

type PayTypeRow struct {
	ID        int       `csv:"ID"`
	BeginTime time.Time `csv:"BEGIN_TIME" format:"2006-01-02 15:04:05"`
	EndTime   time.Time `csv:"END_TIME"`
	Descr     string    `csv:"DESCRIPTION"`
	RegionID  int       `csv:"REGION_ID"`
}

func (a *PayTypes) Render(db *sqlx.DB) (r []string, err error) {
	var tps = []PayTypeRow{
		{ID: 0, Descr: "Наличные", BeginTime: time.Unix(0, 0).UTC()},
		{ID: 1, Descr: "Банк", BeginTime: time.Unix(0, 0).UTC()},
		{ID: 2, Descr: "Внешние платежи", BeginTime: time.Unix(0, 0).UTC()},
		{ID: 3, Descr: "Credit Card", BeginTime: time.Unix(0, 0).UTC()},
		{ID: 4, Descr: "Бонус", BeginTime: time.Unix(0, 0).UTC()},
		{ID: 5, Descr: "Корректировка", BeginTime: time.Unix(0, 0).UTC()},
		{ID: 6, Descr: "Компенсация", BeginTime: time.Unix(0, 0).UTC()},
		{ID: 7, Descr: "Перевод личных средств", BeginTime: time.Unix(0, 0).UTC()},
		{ID: 8, Descr: "Пересчитать", BeginTime: time.Unix(0, 0).UTC()},
		{ID: 44, Descr: "SberbankNew", BeginTime: time.Unix(0, 0).UTC()},
		{ID: 67, Descr: "Sberbank", BeginTime: time.Unix(0, 0).UTC()},
		{ID: 68, Descr: "РИРЦ", BeginTime: time.Unix(0, 0).UTC()},
		{ID: 101, Descr: "РИРЦ кабельное", BeginTime: time.Unix(0, 0).UTC()},
		{ID: 102, Descr: "РИРЦ Интернет", BeginTime: time.Unix(0, 0).UTC()},
		{ID: 110, Descr: "Сбербанк карты", BeginTime: time.Unix(0, 0).UTC()},
	}

	r = csv.MarshalCSV(tps, ";", "")
	return r, nil
}

func (a *PayTypes) GetFileName() string {
	return fmt.Sprintf("PAY_TYPE_%s.txt", time.Now().Format("20060102_1504"))
}
