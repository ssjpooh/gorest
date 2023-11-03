package roomsapi

import (
	"net/http"

	rooms "restApi/model/vbase"
	dbHandler "restApi/util/db"

	logger "restApi/util/log"

	authHandler "restApi/util/auth"

	"github.com/gin-gonic/gin"
)

func getRoomList(context *gin.Context) []rooms.Rooms {

	var roomList []rooms.Rooms

	query := dbHandler.MakeQuery(dbHandler.SELECT, rooms.RoomsColumns, dbHandler.FROM, "rooms")

	err := dbHandler.Db.Select(&roomList, query)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "not Found"})
		logger.Logger(logger.GetFuncNm(), "select error : ", err.Error())
	}
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
