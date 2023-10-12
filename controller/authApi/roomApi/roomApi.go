package roomapi

import (
	"net/http"
	authHandler "restApi/util/auth"

	rooms "restApi/model/rooms"
	dbHandler "restApi/util/db"

	logger "restApi/util/log"

	"github.com/gin-gonic/gin"
)

func getRoomList(context *gin.Context) []rooms.Room {
	var roomList []rooms.Room

	query := "SELECT ROOM_NO, TITLE FROM ROOM_TBL "

	err := dbHandler.Db.Select(&roomList, query)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "not Found"})
		logger.Logger(logger.GetFuncNm(), "select error : ", err.Error())
	}

	return roomList

}
func RoomApiHandler(v1 *gin.RouterGroup) {

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
}
