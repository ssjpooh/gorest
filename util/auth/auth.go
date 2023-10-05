package auth

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	outhInfo "restApi/model/auth"

	"github.com/dgrijalva/jwt-go"
)

func Authenticate(c *gin.Context) {

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization_header_missing"})
		return
	}

	bearerToken := strings.Trim(authHeader[len("Bearer "):], " ")

	token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
		return outhInfo.JWTKey, nil
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		requested := claims["requested"].(string)
		if requested != c.RemoteIP() {
			log.Println()
			log.Println("requested  : ", requested)
			log.Println("c.RemoteIP() : ", c.RemoteIP())
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid_token"})
			return
		}

	} else {

		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid_token"})
		return
	}

	c.Next()
}
