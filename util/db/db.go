package db

import (
	"log"

	"github.com/jmoiron/sqlx"
)

var Db *sqlx.DB

/*
Description : db 접속
Params      : lastRequestDt
return      :
Author      : ssjpooh
Date        : 2023.10.10
*/
func DbConnect() {

	var err error
	dsn := "root:root@123@tcp(localhost:3306)/foxedu"
	Db, err = sqlx.Open("mysql", dsn)
	if err != nil {
		log.Print(err)
	}
}
