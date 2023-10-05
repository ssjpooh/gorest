package memoryapi

import (
	"log"
	"net/http"

	gmap "restApi/util/memory"

	outhInfo "restApi/model/auth"

	"github.com/gin-gonic/gin"
)

func getGlobalAuth(c *gin.Context, token string) outhInfo.AuthInfo {

	return gmap.GetAuthInfo(token)
}

func MemoryApiHaneler(v1 *gin.RouterGroup) {

	v1.GET("/memory/:token", func(context *gin.Context) {
		token := context.Param("token")
		log.Println("memory token ::: ", token)
		globalMap := getGlobalAuth(context, token)
		log.Println(" GlobalAuthInfoMap ::: ", gmap.GlobalAuthInfoMap)
		context.IndentedJSON(http.StatusOK, globalMap)
	})
}
