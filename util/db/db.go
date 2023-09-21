package db

import (
	"log"

	"github.com/jmoiron/sqlx"
)

var Db *sqlx.DB

func DbConnect() {

	var err error
	dsn := "root:root@123@tcp(localhost:3306)/foxedu"
	Db, err = sqlx.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
}
