package gen

import "time"

type Config struct {
	MySQL       string    `env:"MYSQL"`
	MyNames     string    `env:"NAMES"`
	Path        string    `env:"PATH" default:"./tmp"`
	RegionName  string    `env:"REGION_NAME" default:"Основной"`
	Country     string    `env:"COUNTRY" default:"Россия"`
	RegionID    int       `env:"REGION_ID" default:"1"`
	InitDateS   string    `env:"INIT_DATE" default:"2000-01-01"`
	InitDate    time.Time `env:"-"`
	OnlyOneDay  bool      `env:"ONLY_ONE_DAY"`
	CompanyCode string    `env:"URPREF" default:"720000"`
}

func (c *Config) CalcInitDate() {
	var err error
	c.InitDate, err = time.Parse("2006-01-02", c.InitDateS)
	if c.InitDate.IsZero() || err != nil {
		c.InitDate, _ = time.Parse("2006-01-02", "2000-01-01")
		return
	}
}
