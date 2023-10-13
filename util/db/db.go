package db

import (
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"

	logger "restApi/util/log"
	options "restApi/util/options"
)

var Db *sqlx.DB

var SELECT = "SELECT "
var UPDATE = "UPDATE "
var INSERT = "INSERT INTO "
var DELETE = "DELETE "
var FROM = "FROM "
var WHERE = "WHERE "
var SET = "SET "

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

func ColumnsForStruct(s interface{}) []string {
	var columns []string
	value := reflect.ValueOf(s)
	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)
		if colTag, ok := field.Tag.Lookup("db"); ok {
			columns = append(columns, colTag)
		}
	}
	return columns
}

func MakeQuery(partOfQuery ...string) string {

	return strings.Join(partOfQuery, "")
}
