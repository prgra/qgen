package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/prgra/qgen/gen"
)

func main() {
	db, err := sqlx.Connect("mysql" /* driver name */, os.Getenv("QGEN_MYSQL") /* data source name */)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	_, err = db.Exec(fmt.Sprintf("SET NAMES %s", "latin1"))
	if err != nil {
		log.Println(err)
		os.Exit(2)
	}
	err = gen.WriteToFile(&gen.Abons{}, db)
	if err != nil {
		log.Println(err)
		os.Exit(3)
	}
}