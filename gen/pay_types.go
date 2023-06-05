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
		{ID: 0, Descr: "Наличные", BeginTime: EnvInitDate, RegionID: EnvRegionID},
		{ID: 1, Descr: "Банк", BeginTime: EnvInitDate, RegionID: EnvRegionID},
		{ID: 2, Descr: "Внешние платежи", BeginTime: EnvInitDate, RegionID: EnvRegionID},
		{ID: 3, Descr: "Credit Card", BeginTime: EnvInitDate, RegionID: EnvRegionID},
		{ID: 4, Descr: "Бонус", BeginTime: EnvInitDate, RegionID: EnvRegionID},
		{ID: 5, Descr: "Корректировка", BeginTime: EnvInitDate, RegionID: EnvRegionID},
		{ID: 6, Descr: "Компенсация", BeginTime: EnvInitDate, RegionID: EnvRegionID},
		{ID: 7, Descr: "Перевод личных средств", BeginTime: EnvInitDate, RegionID: EnvRegionID},
		{ID: 8, Descr: "Пересчитать", BeginTime: EnvInitDate, RegionID: EnvRegionID},
		{ID: 44, Descr: "SberbankNew", BeginTime: EnvInitDate, RegionID: EnvRegionID},
		{ID: 67, Descr: "Sberbank", BeginTime: EnvInitDate, RegionID: EnvRegionID},
		{ID: 68, Descr: "РИРЦ", BeginTime: EnvInitDate, RegionID: EnvRegionID},
		{ID: 101, Descr: "РИРЦ кабельное", BeginTime: EnvInitDate, RegionID: EnvRegionID},
		{ID: 102, Descr: "РИРЦ Интернет", BeginTime: EnvInitDate, RegionID: EnvRegionID},
		{ID: 110, Descr: "Сбербанк карты", BeginTime: EnvInitDate, RegionID: EnvRegionID},
	}

	r = csv.MarshalCSV(tps, ";", "")
	return r, nil
}

func (a *PayTypes) GetFileName() string {
	return fmt.Sprintf("PAY_TYPE_%s.txt", time.Now().Format("20060102_1504"))
}
