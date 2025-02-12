package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigtoml"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/config"
	"github.com/prgra/qgen/gen"
	"github.com/prgra/qgen/yhnt"
)

type Report struct {
	G       gen.Generator
	ErrCode int
}

func main() {

	var cfg config.Config
	loader := aconfig.LoaderFor(&cfg, aconfig.Config{
		SkipFlags: true,
		EnvPrefix: "YHNT",
		Files:     []string{"/etc/qgen.toml", "qgen.toml"},
		FileDecoders: map[string]aconfig.FileDecoder{
			".toml": aconfigtoml.New(),
		},
	})
	log.SetOutput(os.Stdout)
	if err := loader.Load(); err != nil {
		panic(err)
	}
	cfg.CalcInitDate()

	db, err := sqlx.Connect("mysql", cfg.MySQL)
	if err != nil {
		log.Println("mysql", err)
		os.Exit(1)
	}
	if cfg.MyNames != "" {
		_, err = db.Exec(fmt.Sprintf("SET NAMES %s", cfg.MyNames))
		if err != nil {
			log.Println(err)
			os.Exit(2)
		}
	}
	if cfg.OnlyOneDay {
		dberr := db.Get(&cfg.InitDate, "SELECT date FROM qgenlog ORDER BY id DESC LIMIT 1")
		if dberr != nil {
			fmt.Println("get last date", dberr)
		}
		fmt.Println("init date", cfg.InitDate)
	}
	t := time.Now()
	var reports = []Report{
		{&yhnt.Abons{}, 1},
		{&yhnt.GateWay{}, 2},
		{&yhnt.DocType{}, 3},
		{&yhnt.Payments{}, 4},
		{&yhnt.PaymentsType{}, 5},
		{&yhnt.IPPlan{}, 6},
	}
	for _, r := range reports {
		err = gen.UploadToFTP(r.G, cfg, db)
		if err != nil {
			log.Printf("report %s :%v\n", r.G.GetRemoteDir()+r.G.GetFileName(), err)
			os.Exit(r.ErrCode)
		}
	}
	_, err = db.Exec("INSERT INTO qgenlog (date, comment) values(?,?)", t, "OK")
	if err != nil {
		fmt.Println("can't insert into log", err)
	}

}
