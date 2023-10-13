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

func TokenApiHandler(router *gin.Engine) {

	router.POST("/oauth/token", tokenHandler)
	router.POST("/oauth/refresh", RefreshTokenHandler)
}

/*
Description : refresh token 으로 JWT token 발급
Params      : gin.Context
return      : JSON(token info)
Author      : ssjpooh
Date        : 2023.10.10
*/
func RefreshTokenHandler(c *gin.Context) {
	clientId := ""
	res := true
	refresh := c.PostForm("refresh_token")

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
			res = insertToken(c, token, clientId, expiry, refresh, util.GetLocalIP())
			if !res {
				logger.Logger(logger.GetFuncNm(), "refresh token insert err")
				return
			} else {
				c.JSON(http.StatusOK, gin.H{"access_token": token, "refresh_token": refresh, "token_type": "bearer", "expires_in": expiry})
			}
		}
	}
}

func TokenApi(c *gin.Context, args ...string) {

	var clientID string
	var clientSecret string

	if len(args) > 0 {
		for index, arg := range args {
			if index == 0 {
				clientID = arg
			}
			if index == 1 {
				clientSecret = arg
			}
		}
	} else {
		logger.Logger(logger.GetFuncNm(), "param err : client id , client secret check ")
		return
	}

	res := true
	var client oauthInfo.ClientDetails

	oauthClientDetailsSelect := dbHandler.MakeQuery(dbHandler.SELECT, oauthInfo.OAuthClientDetailsColumns, dbHandler.FROM, "OAUTH_CLIENT_DETAILS ", dbHandler.WHERE, "CLIENT_ID = ? ")
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

	var oauth oauthInfo.OauthInfo
	var expiry = time.Now().Add(time.Second * TokenExpiry).Unix()
	token, refreshToken, serverAddr, err := GenerateToken(c, expiry, clientID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error generator token"})
		return
	}

	milliseconds := expiry

	oauthClientTokensSelect := dbHandler.MakeQuery(dbHandler.SELECT, oauthInfo.OAuthClientTokensColumns, dbHandler.FROM, "OAUTH_CLIENT_TOKENS ", dbHandler.WHERE, "CLIENT_ID = ? ")
	logger.Logger(logger.GetFuncNm(), " SELECT : ", oauthClientTokensSelect, " client id : ", clientID)
	err = dbHandler.Db.Get(&oauth, oauthClientTokensSelect, clientID)
	if err != nil {
		if err == sql.ErrNoRows {
			res = insertToken(c, token, clientID, milliseconds, refreshToken, serverAddr)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error select oauth tokens error "})
			return
		}
	} else {
		getExpiresAt := oauth.ExpiresAT
		nowMileSec := expiry
		// 시간이 지나면 기존의 token 에 delete 하고 새롭게 insert 하여 client id 당 1개의 토큰을 유지한다.
		if nowMileSec > getExpiresAt {
			res = deleteToken(c, clientID)
			if res {
				res = insertToken(c, token, clientID, milliseconds, refreshToken, serverAddr)
			}
		} else {
			// 안지 났다
			token = oauth.Token
			milliseconds = oauth.ExpiresAT
			gmap.GetAuthInfo(token)
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

	TokenApi(c, clientID, clientSecret)
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
func insertToken(c *gin.Context, token string, clientID string, milliseconds int64, refreshToken string, serverAddr string) bool {

	_, err := dbHandler.Db.Exec("INSERT INTO OAUTH_CLIENT_TOKENS (token, client_id, expires_at, refresh_token, server_address) VALUES (?, ?, ?, ?, ? )", token, clientID, milliseconds, refreshToken, serverAddr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error insert tokens"})
		return false
	} else {
		gmap.SetAuthInfo(token, clientID, serverAddr, 0, milliseconds, time.Now().Unix())
		return true
	}

}

func updateToken(c *gin.Context, token, clientId string) bool {
	_, err := dbHandler.Db.Exec("UPDATE OAUTH_CLIENT_TOKENS SET TOKEN = ? WHERE CLIENT_ID = ? ", token, clientId)
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

	_, err := dbHandler.Db.Exec("DELETE FROM OAUTH_CLIENT_TOKENS WHERE CLIENT_ID = ? ", clientID)

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
	_, err := dbHandler.Db.Exec("DELETE FROM OAUTH_CLIENT_TOKENS WHERE REFRESH_TOKEN = ? ", refreshToken)

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
func GenerateToken(c *gin.Context, exp int64, clientId string) (string, string, string, error) {

	token := jwt.New(jwt.SigningMethodHS256)
	serverAddr := util.GetLocalIP()
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["user"] = "ssj"
	claims["exp"] = exp
	claims["requested"] = c.RemoteIP()
	claims["server"] = serverAddr

	tokenString, err := token.SignedString(oauthInfo.JWTKey)
	if err != nil {
		return "", "", "", err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = "ssj"
	rtClaims["clientId"] = clientId
	rtClaims["exp"] = time.Now().Add(time.Second * TokenExpiry * 24).Unix()

	rt, err := refreshToken.SignedString(oauthInfo.JWTKey)
	if err != nil {
		return "", "", "", err
	}

	return tokenString, rt, serverAddr, nil
}
