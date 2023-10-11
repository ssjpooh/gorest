package tokenapi

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	oauthInfo "restApi/model/auth"
	"restApi/util"
	dbHandler "restApi/util/db"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	authHandler "restApi/util/auth"
	gmap "restApi/util/memory"
)

const (
	TokenExpiry = 3600 // 1 hour in seconds
	// TokenExpiry = 180
)

func TokenApiHandler(router *gin.Engine) {

	router.POST("/oauth/token", tokenHandler)

	v1 := router.Group("/v1")
	v1.POST("/oauth/refresh", authHandler.Authenticate, refreshTokenHandler)
}

func refreshTokenHandler(c *gin.Context) {

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization_header_missing"})
		return
	}

	bearerToken := strings.Trim(authHeader[len("Bearer "):], " ")
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
		var oauth oauthInfo.OauthInfo
		err = dbHandler.Db.Get(&oauth, "SELECT refresh_token, client_id, expires_at, token, server_address from oauth_tokens where refresh_token = ? ", refresh)
		if err != nil {
			// 오류
			c.JSON(http.StatusUnauthorized, gin.H{"error": "search refresh token"})
			return
		} else {
			res = updateToken(c, token, oauth.ClientID)
		}

		if res {
			gmap.PatchAuthInfo(bearerToken, token)
		}
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

	res := true
	clientID := c.PostForm("client_id")
	clientSecret := c.PostForm("client_secret")

	var client oauthInfo.ClientDetails

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

	var oauth oauthInfo.OauthInfo
	var expiry = time.Now().Add(time.Second * TokenExpiry).Unix()
	token, refreshToken, serverAddr, err := generateToken(c, expiry)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error generator token"})
		return
	}

	milliseconds := expiry
	err = dbHandler.Db.Get(&oauth, "SELECT client_id, expires_at, token from oauth_tokens where client_id = ? ", clientID)
	if err != nil {
		if err == sql.ErrNoRows {

			gmap.SetAuthInfo(token, clientID, serverAddr, 0, milliseconds, time.Now().Unix())
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error insert tokens"})
		return false
	} else {
		return true
	}

}

func updateToken(c *gin.Context, token, clientId string) bool {
	_, err := dbHandler.Db.Exec("UPDATE OAUTH_TOKENS SET TOKEN = ? WHERE CLIENT_ID = ? ", token, clientId)
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

	_, err := dbHandler.Db.Exec("DELETE FROM OAUTH_TOKENS WHERE CLIENT_ID = ? ", clientID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error delete tokens"})
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
func generateToken(c *gin.Context, exp int64) (string, string, string, error) {

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
	rtClaims["exp"] = time.Now().Add(time.Second * TokenExpiry * 24).Unix()

	rt, err := refreshToken.SignedString(oauthInfo.JWTKey)
	if err != nil {
		return "", "", "", err
	}

	return tokenString, rt, serverAddr, nil
}
