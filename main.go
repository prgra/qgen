package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/gen"
)

type Report struct {
	G       gen.Generator
	ErrCode int
}

func main() {
	db, err := sqlx.Connect("mysql", os.Getenv("QGEN_MYSQL"))
	if err != nil {
		log.Println("mysql", err)
		os.Exit(1)
	}
	if os.Getenv("QGEN_NAMES") != "" {
		_, err = db.Exec(fmt.Sprintf("SET NAMES %s", os.Getenv("QGEN_NAMES")))
		if err != nil {
			log.Println(err)
			os.Exit(2)
		}
	}
	gen.EnvCompanyCode = "720000"
	if os.Getenv("QGEN_URPREF") != "" {
		gen.EnvCompanyCode = os.Getenv("QGEN_URPREF")
	}
	gen.EnvRegionID, _ = strconv.Atoi(os.Getenv("QGEN_REGION_ID"))
	gen.EnvRegionName = os.Getenv("QGEN_REGION_NAME")
	if gen.EnvRegionName == "" {
		gen.EnvRegionName = "Основной"
	}
	gen.EnvOnlyOneDay = os.Getenv("QGEN_ONLY_ONE_DAY") == "1" ||
		strings.ToLower(os.Getenv("QGEN_ONLY_ONE_DAY")) == "true"
	gen.EnvInitDate, err = time.Parse("2006-01-02", os.Getenv("QGEN_INIT_DATE"))
	if err != nil {
		gen.EnvInitDate = time.Unix(0, 0)
	}
	gen.EnvInitDate = gen.EnvInitDate.UTC()
	gen.EnvCountry = os.Getenv("QGEN_COUNTRY")

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
	}
	for _, r := range reports {
		err = gen.WriteToFile(r.G, db)
		if err != nil {
			log.Println(err)
			os.Exit(r.ErrCode)
		}
	}
}
