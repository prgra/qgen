package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/gen"
)

type Report struct {
	G       gen.Generator
	ErrCode int
}

func main() {
	db, err := sqlx.Connect("mysql" /* driver name */, os.Getenv("QGEN_MYSQL") /* data source name */)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	if os.Getenv("QGEN_NAMES") != "" {
		_, err = db.Exec(fmt.Sprintf("SET NAMES %s", os.Getenv("QGEN_NAMES")))
		if err != nil {
			log.Println(err)
			os.Exit(2)
		}
	}
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
	}
	for _, r := range reports {
		err = gen.WriteToFile(r.G, db)
		if err != nil {
			log.Println(err)
			os.Exit(r.ErrCode)
		}
	}
}
