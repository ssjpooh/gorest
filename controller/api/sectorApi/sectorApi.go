package sectorapi

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	sector "restApi/model/vbase"

	dbHandler "restApi/util/db"
	logger "restApi/util/log"
)

func getSectorInfoHandler(context *gin.Context) sector.Sectors {

	var sectorInfo []sector.Sectors

	sectorSelectListQuery := dbHandler.MakeQuery(dbHandler.SELECT, sector.SectorColumns, dbHandler.FROM, "sectors ", dbHandler.Limit(1))

	err := dbHandler.Db.Select(&sectorInfo, sectorSelectListQuery)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "not Found"})
		logger.Logger(logger.GetFuncNm(), "select error ", err.Error())
	}

	return sectorInfo[0]
}

func insertSectorInfoHandler(context *gin.Context) sql.Result {
	var newSectorInfo sector.Sectors
	err := context.ShouldBindJSON(&newSectorInfo)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		logger.Logger(logger.GetFuncNm(), "insert error : ", err.Error())
	}

	query := "UPDATE USER_TBL set kor_user_name = ? , eng_user_name = ? where user_id = ? "
	result, err := dbHandler.Db.Exec(query, newSectorInfo.Sector)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		logger.Logger(logger.GetFuncNm(), "patch error : ", err.Error())
	}

	return result
}

func SectorApiHandler(router *gin.RouterGroup) {

	router.GET("/sectorInfo", func(context *gin.Context) {
		sectorInfo := getSectorInfoHandler(context)
		context.IndentedJSON(http.StatusOK, sectorInfo)

	})
	router.POST("/sectorInfo", func(context *gin.Context) {
		sectorInfo := insertSectorInfoHandler(context)
		context.IndentedJSON(http.StatusOK, sectorInfo)

	})
}
