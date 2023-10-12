package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	tokenapi "restApi/controller/api/tokenApi"
	memberApi "restApi/controller/authApi/memberApi"
	memoryApi "restApi/controller/authApi/memoryApi"
	roomApi "restApi/controller/authApi/roomApi"

	dbHandler "restApi/util/db"

	_ "restApi/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	gmap "restApi/util/memory"

	logger "restApi/util/log"
	options "restApi/util/options"
)

var logFile = time.Now().Format("2006-01-02 15:04")

// @title Swagger Example API
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/
// @BasePath /v1

// @SecurityDefinitions.apikey Authorization
// @type apiKey
// @in header
// @name Authorization
func main() {

	logger.SetFileName(logFile)

	logger.Logger(logger.GetFuncNm(), "main Start")
	// map 설정
	gmap.Init()

	// option 설정
	options.Init()

	dbHandler.DbConnect()
	defer dbHandler.Db.Close()

	router := gin.Default()

	// v1 그룸에 포함 되지 않는 api들
	router.GET("/", indexHandler)
	tokenapi.TokenApiHandler(router)

	v1 := router.Group("/v1")
	memberApi.MmberApiHandler(v1)
	roomApi.RoomApiHandler(v1)
	memoryApi.MemoryApiHaneler(v1)

	url := ginSwagger.URL("https://local.foxedu.kr:443/swagger/doc.json") // The url pointing to API definition
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	err := router.RunTLS(":443", "/Users/shinsangjun/go/src/restApi/util/cert/STAR.foxedu.kr.crt", "/Users/shinsangjun/go/src/restApi/util/cert/STAR.foxedu.kr.key")

	if err != nil {
		panic(err)
	}

}

// description : 접속 확인
func indexHandler(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, "ok")
}
