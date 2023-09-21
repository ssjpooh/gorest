package auth

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
	"time"

	dbHandler "restApi/util/db"

	"github.com/gin-gonic/gin"

	outhInfo "restApi/model/auth"
)

func Authenticate(c *gin.Context) {

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization_header_missing"})
	}

	bearerToken := strings.Trim(authHeader[len("Bearer "):], " ")

	var tokenExpiry outhInfo.OauthInfo
	err := dbHandler.Db.Get(&tokenExpiry, "SELECT expires_at FROM oauth_tokens WHERE token=?", bearerToken)
	if err != nil {
		if err == sql.ErrNoRows {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid_token"})
			return
		}
		log.Fatal(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}

	nowMileSec := time.Now().UnixNano() / int64(time.Millisecond)
	if nowMileSec > tokenExpiry.ExpiresAT {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token_expired"})
		return
	}

	c.Next()
}
