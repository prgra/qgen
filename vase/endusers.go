package vase

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/config"
	"github.com/prgra/qgen/csv"
)

// EndUsers выгрузка связей оконечных пользователей с юрлицом end_users_v1 (объектная).
// Выгружается после subscribers: все абоненты должны быть уже загружены.
type EndUsers struct{}

type EndUserRow struct {
	// ID абонента оконечного пользователя из subscribers.
	EndUserID string `db:"-" csv:"end_user_id"`
	// ID физ лица, для плоского формата subscribers - пустое поле.
	PersonID string `db:"-" csv:"person_id"`
	// ID абонента юридического лица из subscribers.
	CompanyID string `db:"uid" csv:"company_id"`
	// ID юр лица, для плоского формата subscribers - пустое поле.
	ULID string `db:"-" csv:"ul_id"`
}

func (a *EndUsers) Render(db *sqlx.DB, cfg config.Config) (r []string, err error) {
	var rows []EndUserRow
	dta := cfg.InitDate.Format("2006-01-02")
	if cfg.OnlyOneDay {
		dta = time.Now().Format("2006-01-02")
	}
	// те же юрлица, что попадают в subscribers
	err = db.Select(&rows, `select u.uid
from
users u
JOIN dv_main dv ON dv.uid=u.uid
LEFT JOIN admin_actions aa2 on aa2.id = (select id from admin_actions
	where uid=u.uid order by id desc limit 1)
JOIN tarif_plans tp ON tp.id=dv.tp_id
WHERE u.company_id > 0 AND aa2.datetime >= ?
order by u.uid`, dta)
	if err != nil {
		return nil, err
	}
	for i := range rows {
		rows[i].EndUserID = endUserID(rows[i].CompanyID)
	}
	r = csv.MarshalCSVNoHeader(rows, ";", "\"")
	return r, nil
}

func (a *EndUsers) GetFileName() string {
	return fmt.Sprintf("end_users_v1_%s.txt", time.Now().UTC().Format("02012006150405"))
}

func (a *EndUsers) GetRemoteDir() string {
	return ""
}
