package memoryapi

import (
	"net/http"

	gmap "restApi/util/memory"

	outhInfo "restApi/model/auth"

	"github.com/gin-gonic/gin"
)

func getGlobalAuth(c *gin.Context, token string) (outhInfo.AuthInfo, error) {
	return gmap.GetAuthInfo(token, c.Request.RequestURI)
}

func MemoryApiHaneler(v1 *gin.RouterGroup) {

	v1.GET("/memory/:token", func(context *gin.Context) {
		token := context.Param("token")
		globalMap, err := getGlobalAuth(context, token)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid_token"})
			return
		}
		context.IndentedJSON(http.StatusOK, globalMap)
	})
}
