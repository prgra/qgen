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
	"github.com/prgra/qgen/gen"
)

type Report struct {
	G       gen.Generator
	ErrCode int
}

func main() {

	var cfg gen.Config
	loader := aconfig.LoaderFor(&cfg, aconfig.Config{
		SkipFlags: true,
		EnvPrefix: "QGEN",
		Files:     []string{"/etc/qgen.toml", "qgen.toml"},
		FileDecoders: map[string]aconfig.FileDecoder{
			".toml": aconfigtoml.New(),
		},
	})

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
		{&gen.DocType{}, 1},
		{&gen.Abons{}, 2},
		{&gen.AbonIdent{}, 3},
		{&gen.AbonAddr{}, 4},
		{&gen.Region{}, 5},
		{&gen.PayTypes{}, 6},
		{&gen.Supplementary{}, 7},
		{&gen.IPPlan{}, 8},
		{&gen.GateWay{}, 9},
		{&gen.Payments{}, 10},
		{&gen.AbonUsers{}, 11},
		{&gen.AbonSrv{}, 12},
		{&gen.IPGateWay{}, 13},
	}
	for _, r := range reports {
		err = gen.WriteToFile(r.G, cfg, db)
		if err != nil {
			log.Printf("report %s :%v\n", r.G.GetFileName(), err)
			os.Exit(r.ErrCode)
		}
	}
	_, err = db.Exec("INSERT INTO qgenlog (date, comment) values(?,?)", t, "OK")
	if err != nil {
		fmt.Println("can't insert into log", err)
	}
}
