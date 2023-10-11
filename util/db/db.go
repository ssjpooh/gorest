package db

import (
	"log"

	"github.com/jmoiron/sqlx"

	options "restApi/util/options"
)

var Db *sqlx.DB

func DbConnect() {

	var err error
	dsn := options.Prop.Id + ":" + options.Prop.Pw + "@tcp(" + options.Prop.Url + ")/" + options.Prop.Url
	Db, err = sqlx.Open("mysql", dsn)
	if err != nil {
		log.Print(err)
	}
}
