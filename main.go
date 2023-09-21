package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	tokenapi "restApi/controller/api/tokenApi"
	memberApi "restApi/controller/authApi/memberApi"
	roomApi "restApi/controller/authApi/roomApi"
	dbHandler "restApi/util/db"
)

type ClientDetails struct {
	ClientID     string `db:"client_id"`
	ClientSecret string `db:"client_secret"`
}

func main() {

	dbHandler.DbConnect()
	defer dbHandler.Db.Close()

	router := gin.Default()

	// v1 그룸에 포함 되지 않는 api들
	router.GET("/", indexHandler)
	tokenapi.TokenApiHandler(router)

	v1 := router.Group("/v1")
	memberApi.MmberApiHandler(v1)
	roomApi.RoomApiHandler(v1)

	err := router.RunTLS(":443", "/Users/shinsangjun/go/src/restApi/util/cert/STAR.foxedu.kr.crt", "/Users/shinsangjun/go/src/restApi/util/cert/STAR.foxedu.kr.key")

	if err != nil {
		panic(err)
	}

}

// description : 접속 확인
func indexHandler(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, "ok")
}
