package db

import (
	"github.com/jmoiron/sqlx"

	logger "restApi/util/log"
	options "restApi/util/options"
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
	dsn := options.Prop.Id + ":" + options.Prop.Pw + "@tcp(" + options.Prop.Url + ")/" + options.Prop.Name

	logger.Logger(logger.GetFuncNm(), "dsn : ", dsn)
	Db, err = sqlx.Open("mysql", dsn)
	if err != nil {
		logger.Logger(logger.GetFuncNm(), "[err] :", err.Error())
	}

	logger.Logger(logger.GetFuncNm(), "DB Connect Success")
}
