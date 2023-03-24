package gen

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/csv"
)

type PayTypes struct{}

type PayTypeRow struct {
	ID          int    `csv:"ID"`
	BeginTime   string `csv:"BEGIN_TIME"`
	EndTime     string `csv:"END_TIME"`
	Description string `csv:"DESCRIPTION"`
	RegionID    int    `csv:"REGION_ID"`
}

func (a *PayTypes) Render(db *sqlx.DB) (r []string, err error) {
	var tps = []PayTypeRow{
		{ID: 0, Description: "Наличные"},
		{ID: 1, Description: "Банк"},
		{ID: 2, Description: "Внешние платежи"},
		{ID: 3, Description: "Credit Card"},
		{ID: 4, Description: "Бонус"},
		{ID: 5, Description: "Корректировка"},
		{ID: 6, Description: "Компенсация"},
		{ID: 7, Description: "Перевод личных средств"},
		{ID: 8, Description: "Пересчитать"},
		{ID: 44, Description: "SberbankNew"},
		{ID: 67, Description: "Sberbank"},
		{ID: 68, Description: "РИРЦ"},
		{ID: 101, Description: "РИРЦ кабельное"},
		{ID: 102, Description: "РИРЦ Интернет"},
		{ID: 110, Description: "Сбербанк карты"},
	}

	r = csv.MarshalCSV(tps, ";", "")
	return r, nil
}

func (a *PayTypes) GetFileName() string {
	return fmt.Sprintf("PAY_TYPE_%s.txt", time.Now().Format("20060102_1504"))
}
