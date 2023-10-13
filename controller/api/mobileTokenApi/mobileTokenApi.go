package mobiletokenapi

import (
	"database/sql"
	"net/http"
	token "restApi/controller/api/tokenApi"
	logger "restApi/util/log"

	"github.com/gin-gonic/gin"

	db "restApi/util/db"

	oauthInfo "restApi/model/auth"
	members "restApi/model/members"
)

func MobileTokenApiHandler(router *gin.Engine) {

	router.POST("/mobile/oauth/token", tokenHandler)
	router.POST("/mobile/oauth/refresh", refreshTokenHandler)
}

func tokenHandler(c *gin.Context) {

	clientID := c.PostForm("client_id")
	clientSecret := c.PostForm("client_secret")
	userId := c.PostForm("user_id")
	userPw := c.PostForm("user_pw")

	if clientID == "" || clientSecret == "" {
		logger.Logger(logger.GetFuncNm(), " token error : client Id or client Secret is not nill ")
		return
	}

	if userId == "" || userPw == "" {
		token.TokenApi(c, clientID, clientSecret)
	} else {
		logger.Logger(logger.GetFuncNm(), "mobile token with user_id : ", userId)
		// 모바일용으로 id / pw 로 client id / secret 을 받아서 새로 설정 한다.

		var userInfo members.Member
		memberSelect := db.MakeQuery(db.SELECT, members.MemberColumns, db.FROM, "USER_TBL ", db.WHERE, "USER_ID = ? ")

		logger.Logger(logger.GetFuncNm(), " select :  ", memberSelect)
		err := db.Db.Get(&userInfo, memberSelect, userId)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_user id "})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
			logger.Logger(logger.GetFuncNm(), "err : ", err.Error())
			return
		}

		var oauthClientDetails oauthInfo.ClientDetails
		oauthClientDetailsSelect := db.MakeQuery(db.SELECT, oauthInfo.OAuthClientDetailsColumns, db.FROM, " OAUTH_CLIENT_DETAILS ", db.WHERE, " OWNER_IDX = ? ")

		logger.Logger(logger.GetFuncNm(), " select :  ", oauthClientDetailsSelect, " owner_idx : ", userInfo.OWNERIDX)

		err = db.Db.Get(&oauthClientDetails, oauthClientDetailsSelect, userInfo.OWNERIDX)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_owner_idx "})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
			logger.Logger(logger.GetFuncNm(), "err : ", err.Error())
			return
		}

		// 임의의 POST 데이터 생성
		token.TokenApi(c, oauthClientDetails.ClientID, oauthClientDetails.ClientSecret)
	}

}

func refreshTokenHandler(c *gin.Context) {

}
