package tokenapi

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	outhInfo "restApi/model/auth"
	dbHandler "restApi/util/db"
)

const (
	TokenExpiry = 3600 // 1 hour in seconds
)

func TokenApiHandler(route *gin.Engine) {

	route.POST("/oauth/token", tokenHandler)
}

// description : token 발급
func tokenHandler(c *gin.Context) {

	clientID := c.PostForm("client_id")
	clientSecret := c.PostForm("client_secret")

	fmt.Printf("[%s] [%s]\n", clientID, clientSecret)
	var client outhInfo.ClientDetails
	err := dbHandler.Db.Get(&client, "SELECT client_id, client_secret FROM oauth_client_details WHERE client_id=?", clientID)
	if err != nil {
		log.Fatal(err)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_client"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}

	fmt.Println("aa")
	if client.ClientSecret != clientSecret {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_client2"})
		return
	}

	var oauth outhInfo.OauthInfo

	token := generateToken()
	expiry := time.Now().Add(time.Second * TokenExpiry)
	milliseconds := expiry.UnixNano() / int64(time.Millisecond)
	err = dbHandler.Db.Get(&oauth, "SELECT client_id, expires_at, token from oauth_tokens where client_id = ? ", clientID)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err = dbHandler.Db.Exec("INSERT INTO oauth_tokens (token, client_id, expires_at) VALUES (?, ?, ?)", token, clientID, milliseconds)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error2"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error1"})
			return
		}
	} else {
		getExpiresAt := oauth.ExpiresAT
		nowMileSec := time.Now().UnixNano() / int64(time.Millisecond)
		// 시간이 지났다
		fmt.Printf(" token : [%s], expires_milisec :  [%d], clientId : [%s] ", token, milliseconds, clientID)
		fmt.Println()
		fmt.Printf("현재 : [%d], DB값 : [%d] ", nowMileSec, getExpiresAt)
		if nowMileSec > getExpiresAt {
			_, err = dbHandler.Db.Exec("UPDATE oauth_tokens SET token = ? , expires_at = ? where client_id = ? ", token, milliseconds, clientID)
			if err != nil {
				log.Fatal(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error3"})
				return
			}
		} else {
			// 안지 났다
			token = oauth.Token
			milliseconds = oauth.ExpiresAT
		}
	}

	c.JSON(http.StatusOK, gin.H{"access_token": token, "token_type": "bearer", "expires_in": milliseconds})
}

func generateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
