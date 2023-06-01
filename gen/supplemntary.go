package gen

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/csv"
)

type SupplementaryRow struct {
	ID          int       `csv:"ID"`                                    // 	Идентификатор услуги
	Mnemonic    string    `csv:"MNEMONIC"`                              //		Мнемоническое обозначение
	BeginTime   time.Time `csv:"BEGIN_TIME" time:"2006-01-02 15:04:05"` // 	Время начала действия
	EndTime     time.Time `csv:"END_TIME" time:"2006-01-02 15:04:05"`   // 	Время конца действия
	Description string    `csv:"DESCRIPTION"`                           // 	Описание
	RegionID    int       `csv:"REGION_ID"`                             // 	Идентификатор оператора связи или структурного подразделения (ссылка на справочник операторов или филиалов)
}

type Supplementary struct{}

func (a *Supplementary) Render(db *sqlx.DB) (r []string, err error) {
	supps := []SupplementaryRow{
		{
			ID:          1,
			Mnemonic:    "INET",
			BeginTime:   time.Unix(0, 0).UTC(),
			Description: "Доступ в Интернет по технологии FTTx",
		},
	}

	r = csv.MarshalCSV(supps, ";", "")
	return r, nil
}

func (a *Supplementary) GetFileName() string {
	return fmt.Sprintf("SUPPLEMENTARY_SERVICE_%s.txt", time.Now().Format("20060102_1504"))
}
