package vase

import (
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/config"
)

// Binds бинд-файл соотношения логинов и IP адресов абонентов binds.txt.
// Нужен, только если абоненты авторизуются НЕ через RADIUS/DIAMETER/TACACS+,
// выгружать не реже раза в 15 минут, в файле только актуальные данные.
type Binds struct{}

type bindRow struct {
	Login string `db:"login"`
	IP    string `db:"ip"`
}

func (a *Binds) Render(db *sqlx.DB, cfg config.Config) (r []string, err error) {
	_ = cfg
	var rows []bindRow
	err = db.Select(&rows, `select u.id as login, INET_NTOA(dv.ip) as ip
		from users u
		JOIN dv_main dv ON dv.uid=u.uid
		WHERE dv.ip <> 0 AND u.deleted = 0 AND u.disable = 0
		order by u.id`)
	if err != nil {
		return nil, err
	}
	// формат строки: Логин;ip[,ip...]; каждому IP - не более одного логина
	ips := make(map[string][]string)
	owner := make(map[string]string)
	var logins []string
	for _, b := range rows {
		login := cleanNewlines(b.Login)
		if login == "" || b.IP == "" || owner[b.IP] != "" {
			continue
		}
		owner[b.IP] = login
		if _, ok := ips[login]; !ok {
			logins = append(logins, login)
		}
		ips[login] = append(ips[login], b.IP)
	}
	sort.Strings(logins)
	for _, l := range logins {
		r = append(r, l+";"+strings.Join(ips[l], ","))
	}
	return r, nil
}

func (a *Binds) GetFileName() string {
	return "binds.txt"
}

func (a *Binds) GetRemoteDir() string {
	return ""
}
