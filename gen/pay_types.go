package gen

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/csv"
)

type PayTypes struct{}

type PayTypeRow struct {
	ID        int    `csv:"ID"`
	BeginTime string `csv:"BEGIN_TIME"`
	EndTime   string `csv:"END_TIME"`
	Descr     string `csv:"DESCRIPTION"`
	RegionID  int    `csv:"REGION_ID"`
}

func (a *PayTypes) Render(db *sqlx.DB) (r []string, err error) {
	var tps = []PayTypeRow{
		{ID: 0, Descr: "Наличные"},
		{ID: 1, Descr: "Банк"},
		{ID: 2, Descr: "Внешние платежи"},
		{ID: 3, Descr: "Credit Card"},
		{ID: 4, Descr: "Бонус"},
		{ID: 5, Descr: "Корректировка"},
		{ID: 6, Descr: "Компенсация"},
		{ID: 7, Descr: "Перевод личных средств"},
		{ID: 8, Descr: "Пересчитать"},
		{ID: 44, Descr: "SberbankNew"},
		{ID: 67, Descr: "Sberbank"},
		{ID: 68, Descr: "РИРЦ"},
		{ID: 101, Descr: "РИРЦ кабельное"},
		{ID: 102, Descr: "РИРЦ Интернет"},
		{ID: 110, Descr: "Сбербанк карты"},
	}

	r = csv.MarshalCSV(tps, ";", "")
	return r, nil
}

func (a *PayTypes) GetFileName() string {
	return fmt.Sprintf("PAY_TYPE_%s.txt", time.Now().Format("20060102_1504"))
}
