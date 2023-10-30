package roomapi

import (
	"net/http"

	rooms "restApi/model/rooms"
	dbHandler "restApi/util/db"

	logger "restApi/util/log"

	authHandler "restApi/util/auth"

	"github.com/gin-gonic/gin"
)

func getRoomList(context *gin.Context) []rooms.Room {

	logger.Logger(logger.GetFuncNm(), "1")
	var roomList []rooms.Room

	query := "SELECT room_no, title FROM CONF_TBL "
	logger.Logger(logger.GetFuncNm(), "2")
	err := dbHandler.Db.Select(&roomList, query)
	if err != nil {
		logger.Logger(logger.GetFuncNm(), "3")
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "not Found"})
		logger.Logger(logger.GetFuncNm(), "select error : ", err.Error())
	}
	logger.Logger(logger.GetFuncNm(), "4")
	return roomList

}
func RoomApiHandler(v1 *gin.RouterGroup) {

	logger.Logger(logger.GetFuncNm(), "start")
	// @Summary get room List
	// @Description Get room List
	// @Tags rooms
	// @Accept  json
	// @Produce  json
	// @Param   id path string true "ID"
	// @Success 200 {object} rooms.Room
	// @Router /rooms [get]
	v1.GET("/rooms", authHandler.Authenticate, func(context *gin.Context) {
		roomList := getRoomList(context)
		context.IndentedJSON(http.StatusOK, roomList)
	})

	v1.GET("/roomtest", func(context *gin.Context) {
		roomList := getRoomList(context)
		context.IndentedJSON(http.StatusOK, roomList)
	})

}
