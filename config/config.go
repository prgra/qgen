package config

import "time"

type Config struct {
	MySQL       string    `env:"MYSQL" toml:"mysql"`
	MyNames     string    `env:"NAMES" toml:"mynames"`
	Path        string    `env:"PATH" default:"./tmp" toml:"path"`
	RegionName  string    `env:"REGION_NAME" default:"Основной" toml:"region_name"`
	Country     string    `env:"COUNTRY" default:"Россия" toml:"country"`
	RegionID    int       `env:"REGION_ID" default:"1" toml:"region_id"`
	InitDateS   string    `env:"INIT_DATE" default:"2000-01-01" toml:"init_date"`
	InitDate    time.Time `env:"-"`
	OnlyOneDay  bool      `env:"ONLY_ONE_DAY" default:"false" toml:"only_one_day"`
	CompanyCode string    `env:"URPREF" default:"720000" toml:"company_code"`
	NasGroupID  int       `env:"NAS_GROUP_ID" default:"1" toml:"nas_group_id"`
	CSVChatset  string    `env:"CSV_CHARSET" default:"" toml:"csv_charset"`
	FTPHost     string    `env:"FTP_HOST" default:"" toml:"ftp_host"`
	FTPUser     string    `env:"FTP_USER" default:"" toml:"ftp_user"`
	FTPPass     string    `env:"FTP_PASS" default:"" toml:"ftp_pass"`
}

func (c *Config) CalcInitDate() {
	var err error
	c.InitDate, err = time.Parse("2006-01-02", c.InitDateS)
	if c.InitDate.IsZero() || err != nil {
		c.InitDate, _ = time.Parse("2006-01-02", "2000-01-01")
		return
	}
}
