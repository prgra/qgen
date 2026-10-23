package vase

import (
	"database/sql"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/config"
	"github.com/prgra/qgen/csv"
)

// IPNets выгрузка статических IP адресов и подсетей абонентов ip_nets_v2 (объектная).
type IPNets struct{}

type IPNetRow struct {
	// ID точки подключения, зарезервировано на будущее - пустое поле.
	PointID string `db:"-" csv:"point_id"`
	// ID абонента, соответствует полю id выгрузки subscribers.
	ID string `db:"uid" csv:"id"`
	// IP адрес без маски или подсеть с маской, например 85.202.1.0/30.
	IP string `db:"ip" csv:"ip"`
	// Тип адреса: 0 - серый, 1 - публичный, 3 - VPN.
	IPType int `db:"-" csv:"ip_type"`
	// Дата выдачи адреса.
	BeginDate time.Time `db:"-" csv:"begin_date" time:"02.01.2006T15:04:05"`
	// Дата окончания.
	EndDate time.Time `db:"-" csv:"end_date" time:"02.01.2006T15:04:05"`
	// MAC-адрес в формате 00:00:00:00:00:00, для подсетей не выгружается.
	MAC sql.NullString `db:"mac" csv:"mac"`
	// ID физ лица, для плоского формата subscribers - пустое поле.
	PersonID string `db:"-" csv:"person_id"`

	// Служебные поля выборки, в CSV не выгружаются.
	Mask         sql.NullString `db:"mask" csv:"-"`
	Status       int            `db:"status" csv:"-"`
	ContractDate sql.NullTime   `db:"contract_date" csv:"-"`
	FirstAction  sql.NullTime   `db:"actual_from" csv:"-"`
	LastAction   sql.NullTime   `db:"actual_to" csv:"-"`
}

func (a *IPNets) Render(db *sqlx.DB, cfg config.Config) (r []string, err error) {
	var rows []IPNetRow
	dta := cfg.InitDate.Format("2006-01-02")
	if cfg.OnlyOneDay {
		dta = time.Now().Format("2006-01-02")
	}
	// тот же набор абонентов, что и в subscribers
	err = db.Select(&rows, `select
u.uid,
INET_NTOA(dv.ip) as ip,
INET_NTOA(dv.netmask) as mask,
dh.mac,
u.deleted + u.disable as status,
pi.contract_date as contract_date,
aa1.datetime as actual_from,
aa2.datetime as actual_to
from
users u
JOIN dv_main dv ON dv.uid=u.uid
LEFT JOIN admin_actions aa1 on aa1.id = (select id from admin_actions
	where uid=u.uid order by id limit 1)
LEFT JOIN admin_actions aa2 on aa2.id = (select id from admin_actions
	where uid=u.uid order by id desc limit 1)
LEFT JOIN users_pi pi ON pi.uid=u.uid
LEFT JOIN dhcphosts_hosts dh ON dh.uid=u.uid
WHERE dv.ip <> 0 AND aa2.datetime >= ?`, dta)
	if err != nil {
		return nil, err
	}
	// уникальность записи - ID абонента + IP адрес
	seen := make(map[string]bool)
	var res []IPNetRow
	for i := range rows {
		rows[i].Calc()
		key := rows[i].ID + ";" + rows[i].IP
		if rows[i].IP == "" || seen[key] {
			continue
		}
		seen[key] = true
		res = append(res, rows[i])
	}
	r = csv.MarshalCSVNoHeader(res, ";", "\"")
	return r, nil
}

func (a *IPNetRow) Calc() {
	ip := net.ParseIP(a.IP).To4()
	if ip == nil {
		a.IP = ""
		return
	}
	a.IPType = ipType(ip)
	// подсеть выгружается с маской и без MAC-адреса
	if ones, bits := net.IPMask(net.ParseIP(a.Mask.String).To4()).Size(); bits == 32 && ones > 0 && ones < 32 {
		a.IP = fmt.Sprintf("%s/%d", ip.Mask(net.CIDRMask(ones, 32)), ones)
		a.MAC = sql.NullString{}
	} else {
		a.IP = ip.String()
		a.MAC.String = makeColonMac(a.MAC.String)
	}
	switch {
	case a.ContractDate.Valid && !a.ContractDate.Time.IsZero():
		a.BeginDate = a.ContractDate.Time
	case a.FirstAction.Valid:
		a.BeginDate = a.FirstAction.Time
	}
	if a.Status != 0 && a.LastAction.Valid {
		a.EndDate = a.LastAction.Time
	}
}

var grayNets = []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "100.64.0.0/10"}

// ipType возвращает 0 для серых адресов (RFC1918, CGNAT) и 1 для публичных.
func ipType(ip net.IP) int {
	for _, s := range grayNets {
		_, n, _ := net.ParseCIDR(s)
		if n.Contains(ip) {
			return 0
		}
	}
	return 1
}

// makeColonMac приводит MAC к формату 00:00:00:00:00:00.
func makeColonMac(s string) string {
	m, err := net.ParseMAC(strings.TrimSpace(s))
	if err != nil {
		return ""
	}
	return strings.ToUpper(m.String())
}

func (a *IPNets) GetFileName() string {
	return fmt.Sprintf("ip_nets_v2_%s.txt", time.Now().UTC().Format("02012006150405"))
}

func (a *IPNets) GetRemoteDir() string {
	return ""
}
