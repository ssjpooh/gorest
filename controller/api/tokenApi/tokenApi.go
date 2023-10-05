package tokenapi

import (
	"database/sql"
	"log"
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
		log.Println("111")
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
		log.Println("222")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}

	log.Println("serverAddr : ", serverAddr)
	log.Println("serverAddr : ", len(serverAddr))

	milliseconds := expiry
	err = dbHandler.Db.Get(&oauth, "SELECT client_id, expires_at, token from oauth_tokens where client_id = ? ", clientID)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err = dbHandler.Db.Exec("INSERT INTO oauth_tokens (token, client_id, expires_at, refresh_token, server_address) VALUES (?, ?, ?, ?, ? )", token, clientID, milliseconds, refreshToken, serverAddr)
			if err != nil {
				log.Println("333")
				c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
				return
			}
		} else {
			log.Println("444")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
			return
		}
	} else {
		getExpiresAt := oauth.ExpiresAT
		nowMileSec := expiry
		// 시간이 지났다

		if nowMileSec > getExpiresAt {
			_, err = dbHandler.Db.Exec("UPDATE oauth_tokens SET token = ? , expires_at = ? where client_id = ? ", token, milliseconds, clientID)
			if err != nil {
				log.Println("555")
				c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
				return
			}
		} else {
			// 안지 났다
			token = oauth.Token
			milliseconds = oauth.ExpiresAT
		}
	}

	c.JSON(http.StatusOK, gin.H{"access_token": token, "refresh_token": refreshToken, "token_type": "bearer", "expires_in": milliseconds})
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
