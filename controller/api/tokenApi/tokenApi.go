package tokenapi

import (
	"database/sql"
	"net/http"
	"time"

	outhInfo "restApi/model/auth"
	"restApi/util"
	dbHandler "restApi/util/db"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const (
	TokenExpiry = 3600 // 1 hour in seconds
)

var expiry = time.Now().Add(time.Second * TokenExpiry).Unix()

func TokenApiHandler(route *gin.Engine) {

	route.POST("/oauth/token", tokenHandler)
}

// description : token 발급
func tokenHandler(c *gin.Context) {

	clientID := c.PostForm("client_id")
	clientSecret := c.PostForm("client_secret")

	var client outhInfo.ClientDetails
	err := dbHandler.Db.Get(&client, "SELECT client_id, client_secret FROM oauth_client_details WHERE client_id=?", clientID)
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

	var oauth outhInfo.OauthInfo

	token, refreshToken, serverAddr, err := generateToken(c, expiry)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}

	milliseconds := expiry
	err = dbHandler.Db.Get(&oauth, "SELECT client_id, expires_at, token from oauth_tokens where client_id = ? ", clientID)
	if err != nil {
		if err == sql.ErrNoRows {
			// _, err = dbHandler.Db.Exec("INSERT INTO oauth_tokens (token, client_id, expires_at, refresh_token, server_address) VALUES (?, ?, ?, ?, ? )", token, clientID, milliseconds, refreshToken, serverAddr)
			// if err != nil {
			// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
			// 	return
			// }
			insertToken(c, token, clientID, milliseconds, refreshToken, serverAddr)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
			return
		}
	} else {
		getExpiresAt := oauth.ExpiresAT
		nowMileSec := expiry
		// 시간이 지나면 기존의 token 에 delete 하고 새롭게 insert 하여 client id 당 1개의 토큰을 유지한다.
		if nowMileSec > getExpiresAt {

			res := deleteToken(c, clientID)
			if res {
				insertToken(c, token, clientID, milliseconds, refreshToken, serverAddr)
			}
			// _, err = dbHandler.Db.Exec("UPDATE oauth_tokens SET token = ? , expires_at = ? where client_id = ? ", token, milliseconds, clientID)
			// if err != nil {
			// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
			// 	return
			// }
		} else {
			// 안지 났다
			token = oauth.Token
			milliseconds = oauth.ExpiresAT
		}
	}

	c.JSON(http.StatusOK, gin.H{"access_token": token, "refresh_token": refreshToken, "token_type": "bearer", "expires_in": milliseconds})
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

	_, err := dbHandler.Db.Exec("INSERT INTO oauth_tokens (token, client_id, expires_at, refresh_token, server_address) VALUES (?, ?, ?, ?, ? )", token, clientID, milliseconds, refreshToken, serverAddr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
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

	_, err := dbHandler.Db.Exec("DELETE FROM OAUTH_TOKENS WHERE CLIENT_ID = ? ", clientID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return false
	} else {
		return true
	}
}

func generateToken(c *gin.Context, exp int64) (string, string, string, error) {

	token := jwt.New(jwt.SigningMethodHS256)
	serverAddr := util.GetLocalIP()
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["user"] = "ssj"
	claims["exp"] = exp
	claims["requested"] = c.RemoteIP()
	claims["server"] = serverAddr

	tokenString, err := token.SignedString(outhInfo.JWTKey)
	if err != nil {
		return "", "", "", err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = "ssj"
	rtClaims["exp"] = time.Now().Add(time.Second * TokenExpiry * 24).Unix()

	rt, err := refreshToken.SignedString(outhInfo.JWTKey)
	if err != nil {
		return "", "", "", err
	}

	return tokenString, rt, serverAddr, nil
}
