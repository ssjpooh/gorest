package tokenapi

import (
	"database/sql"
	"net/http"
	"time"

	oauthInfo "restApi/model/auth"
	"restApi/util"
	dbHandler "restApi/util/db"
	logger "restApi/util/log"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	gmap "restApi/util/memory"
)

const (
	TokenExpiry = 3600 // 1 hour in seconds
	// TokenExpiry = 180
)

func TokenApiHandler(router *gin.RouterGroup) {

	router.POST("/oauth/token", tokenHandler)
	router.POST("/oauth/refresh", refreshTokenHandler)
}

func GetOauthClientDetails(context *gin.Context, userIdx string) oauthInfo.OAuthClientDetails {

	var oauthClientDetails oauthInfo.OAuthClientDetails
	query := dbHandler.MakeQuery(dbHandler.SELECT, oauthInfo.OAuthClientDetailsColumns, dbHandler.FROM, "oauth_client_details ", dbHandler.WHERE, "site_idx = ? ")
	logger.Logger(logger.GetFuncNm(), " SELECT : ", query, " client id : ", userIdx)
	err := dbHandler.Db.Get(&oauthClientDetails, query, userIdx)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "not Found"})
		logger.Logger(logger.GetFuncNm(), "select error :", err.Error())
		return oauthClientDetails
	}

	return oauthClientDetails
}

/*
Description : global refresh token handler
Params      : gin.Context
return      : token json
Author      : ssjpooh
Date        : 2023.10.13
*/

func RefreshTokenGlobalApi(c *gin.Context, refresh string) {
	clientId := ""
	res := true
	token, err := jwt.Parse(refresh, func(token *jwt.Token) (interface{}, error) {
		return oauthInfo.JWTKey, nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	var expiry = time.Now().Add(time.Second * TokenExpiry).Unix()
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		clientId = claims["clientId"].(string)
		// Create a new access token using the same claims
		newToken := jwt.New(jwt.SigningMethodHS256)
		newClaims := newToken.Claims.(jwt.MapClaims)
		for key, val := range claims {
			newClaims[key] = val
		}
		newClaims["exp"] = expiry

		token, err := newToken.SignedString(oauthInfo.JWTKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token signeture"})
			return
		}

		// newToken 을 넣자
		res = deleteRefreshToken(c, refresh)
		if res {
			res = insertToken(c, token, clientId, expiry, refresh, sql.NullString{String: util.GetLocalIP(), Valid: true})
			if !res {
				logger.Logger(logger.GetFuncNm(), "refresh token insert err")
				return
			} else {
				c.JSON(http.StatusOK, gin.H{"access_token": token, "refresh_token": refresh, "token_type": "bearer", "expires_in": expiry})
			}
		}
	}
}

/*
Description : refresh token 으로 JWT token 발급
Params      : gin.Context
return      : JSON(token info)
Author      : ssjpooh
Date        : 2023.10.10
*/
func refreshTokenHandler(c *gin.Context) {

	refresh := c.PostForm("refresh_token")
	RefreshTokenGlobalApi(c, refresh)

}

/*
Description : global token handler
Params      : gin.Context
return      : token json
Author      : ssjpooh
Date        : 2023.10.13
*/
func TokenGlobalApi(c *gin.Context, args ...string) {

	var clientID string
	var clientSecret string
	var userSID string

	if len(args) > 0 {
		for index, arg := range args {
			if index == 0 {
				clientID = arg
			}
			if index == 1 {
				clientSecret = arg
			}
			if index == 2 {
				userSID = arg
			}
		}
	} else {
		logger.Logger(logger.GetFuncNm(), "param err : client id , client secret check ")
		return
	}

	res := true
	var client oauthInfo.OAuthClientDetails

	oauthClientDetailsSelect := dbHandler.MakeQuery(dbHandler.SELECT, oauthInfo.OAuthClientDetailsColumns, dbHandler.FROM, "oauth_client_details ", dbHandler.WHERE, "CLIENT_ID = ? ")
	err := dbHandler.Db.Get(&client, oauthClientDetailsSelect, clientID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_client"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}

	if client.ClientSecret != clientSecret {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_client"})
		return
	}

	var oauth oauthInfo.OAuthClientTokens
	var expiry = time.Now().Add(time.Second * TokenExpiry).Unix()
	token, refreshToken, serverAddr, err := GenerateToken(c, expiry, clientID, userSID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error generator token"})
		return
	}

	milliseconds := expiry

	oauthClientTokensSelect := dbHandler.MakeQuery(dbHandler.SELECT, oauthInfo.OAuthClientTokensColumns, dbHandler.FROM, "oauth_client_tokens ", dbHandler.WHERE, "CLIENT_ID = ? ")
	logger.Logger(logger.GetFuncNm(), " SELECT : ", oauthClientTokensSelect, " client id : ", clientID)
	err = dbHandler.Db.Get(&oauth, oauthClientTokensSelect, clientID)
	if err != nil {
		if err == sql.ErrNoRows {
			// logger.Logger(logger.GetFuncNm(), " select : no rows token start Innsert ")
			res = insertToken(c, token, clientID, milliseconds, refreshToken, serverAddr)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error select oauth tokens error "})
			return
		}
	} else {
		// logger.Logger(logger.GetFuncNm(), "exist rows ")
		getExpiresAt := oauth.ExpiresAt
		nowMileSec := time.Now().Unix()
		// 시간이 지나면 기존의 token 에 delete 하고 새롭게 insert 하여 client id 당 1개의 토큰을 유지한다.
		if nowMileSec > getExpiresAt {
			// logger.Logger(logger.GetFuncNm(), "new is bigger than auth info delete and insert  now : ", fmt.Sprintf(" nowMileSec : %d , auth expire : %d ", nowMileSec, getExpiresAt))
			res = deleteToken(c, clientID)
			if res {
				res = insertToken(c, token, clientID, milliseconds, refreshToken, serverAddr)
			}
		} else {
			// logger.Logger(logger.GetFuncNm(), "auth is bigger than now can use token expire : ", fmt.Sprintf(" nowMileSec : %d , auth expire : %d ", nowMileSec, getExpiresAt))
			// 안 지났다
			token = oauth.Token
			milliseconds = oauth.ExpiresAt
			gmap.GetAuthInfo(token, c.Request.RequestURI)
		}
	}

	if res {
		c.JSON(http.StatusOK, gin.H{"access_token": token, "refresh_token": refreshToken, "token_type": "bearer", "expires_in": milliseconds})
	}
}

/*
Description : JWT token 발급
Params      : gin.Context
return      : JSON(token info)
Author      : ssjpooh
Date        : 2023.10.10
*/
func tokenHandler(c *gin.Context) {

	clientID := c.PostForm("client_id")
	clientSecret := c.PostForm("client_secret")

	TokenGlobalApi(c, clientID, clientSecret)
}

/*
Description : 토큰 등록
Params      : gin.Context
Params      : token
Params      : clientId
Params      : expired date
Params      : serverAddr
return      : bool
Author      : ssjpooh
Date        : 2023-10-26
*/
func insertToken(c *gin.Context, token string, clientID string, milliseconds int64, refreshToken string, serverAddr sql.NullString) bool {

	_, err := dbHandler.Db.Exec("INSERT INTO oauth_client_tokens (token, client_id, expires_at, refresh_token, server_address) VALUES (?, ?, ?, ?, ? )", token, clientID, milliseconds, refreshToken, serverAddr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error insert tokens"})
		return false
	} else {
		gmap.SetAuthInfo(token, clientID, serverAddr, 0, milliseconds, time.Now().Unix(), c.Request.RequestURI)
		return true
	}

}

/*
Description : 토큰 수정
Params      : gin.Context
Params      : token
Params      : clientId
return      : bool
Author      : ssjpooh
Date        : 2023-10-26
*/
func updateToken(c *gin.Context, token, clientId string) bool {
	_, err := dbHandler.Db.Exec("UPDATE oauth_client_tokens SET TOKEN = ? WHERE CLIENT_ID = ? ", token, clientId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error update tokens"})
		return false
	} else {
		return true
	}
}

/*
Description : 토큰 삭제
Params      : gin.Context
Params      : clientId
return      : bool
Author      : ssjpooh
Date        : 2023-10-26
*/
func deleteToken(c *gin.Context, clientID string) bool {

	_, err := dbHandler.Db.Exec("DELETE FROM oauth_client_tokens WHERE CLIENT_ID = ? ", clientID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error delete tokens"})
		return false
	} else {
		return true
	}
}

/*
Description : 리프레시 토큰 삭제
Params      : gin.Context
Params      : refresh token
return      : bool
Author      : ssjpooh
Date        : 2023-10-12
*/
func deleteRefreshToken(c *gin.Context, refreshToken string) bool {
	_, err := dbHandler.Db.Exec("DELETE FROM oauth_client_tokens WHERE REFRESH_TOKEN = ? ", refreshToken)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error delete refresh tokens"})
		return false
	} else {
		return true
	}
}

/*
Description : JWT token 생성
Params      : gin.Context
Params      : expire Date
return      : token
return      : refresh token
return      : add server address
return      : error
Author      : ssjpooh
Date        : 2023.10.10
*/
func GenerateToken(c *gin.Context, exp int64, clientId, userID string) (string, string, sql.NullString, error) {

	token := jwt.New(jwt.SigningMethodHS256)
	serverAddr := util.GetLocalIP()
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["user"] = userID
	claims["exp"] = exp
	claims["requested"] = c.RemoteIP()
	claims["server"] = serverAddr

	tokenString, err := token.SignedString(oauthInfo.JWTKey)
	if err != nil {
		return "", "", sql.NullString{}, err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = userID
	rtClaims["clientId"] = clientId
	rtClaims["exp"] = time.Now().Add(time.Second * TokenExpiry * 24).Unix()

	rt, err := refreshToken.SignedString(oauthInfo.JWTKey)
	if err != nil {
		return "", "", sql.NullString{}, err
	}

	return tokenString, rt, sql.NullString{String: serverAddr, Valid: true}, nil
}
