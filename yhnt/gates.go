package yhnt

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/config"
	"github.com/prgra/qgen/csv"
)

type GateWay struct{}

type GateWayRow struct {
	DepID     int       `db:"-"`                                // идентификатор филиала (число)
	IP        string    `db:"ip"`                               // IP-адрес шлюза (либо V4 либо v6) (строка);
	StartDate time.Time `db:"sdate" time:"02.01.2006 15:04:05"` // дата начала действия шлюза (дата);
	EndDate   time.Time `db:"edate" time:"02.01.2006 15:04:05"` // дата завершения действия шлюза (в случае работы оборудования по настоящее время - пустое значение);
	Descr     string    `db:"name"`                             // описание шлюза (строка);
	// адрес установки шлюза:
	Country  string `db:"country"`  // страна;
	Region   string `db:"region"`   // область (строка);
	District string `db:"district"` // район (строка);
	City     string `db:"city"`     // город/поселок/деревня/аул (строка);
	Street   string `db:"street" `  // улица (строка);
	House    string `db:"house"`    // номер дома, строения (строка);
	// Type :: Типы шлюзов:
	// для SGSN – 0;
	// для GGSN – 1;
	// для SMSC – 2;
	// для GMSC – 3;
	// для HSS – 4;
	// для ТФоП-шлюза – 5;
	// для VoIP-шлюза – 6;
	// для ААА – 7;
	// для NAT – 8.
	Type int `db:"-"` // тип шлюза (число):
}

func (gw *GateWay) Render(db *sqlx.DB, cfg config.Config) (r []string, err error) { //
	var gws []GateWayRow //
	err = db.Select(&gws, `select n.ip, n.name,d.country, d.city,
	d.name as district, s.name as street, b.number as house
	from nas n
	left join builds b on b.id = n.location_id
	left join streets s on s.id = b.street_id
	left join districts d on d.id = s.district_id
	where n.gid = ? order by n.id`, cfg.NasGroupID)
	if err != nil {
		return nil, err
	}

	for i := range gws {
		gws[i].DepID = cfg.RegionID
		gws[i].StartDate = cfg.InitDate
		gws[i].Type = 8
		gws[i].Country = cfg.Country
	}

	r = csv.MarshalCSVNoHeader(gws, ";", `"`)
	return r, nil
}

func (gw *GateWay) GetFileName() string {
	return "gates.csv"
}

func (a *GateWay) GetRemoteDir() string {
	return "dictionaries/gates"
}
