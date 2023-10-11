package auth

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	outhInfo "restApi/model/auth"
	memory "restApi/util/memory"

	"github.com/dgrijalva/jwt-go"
)

/*
Description : auth 인증
Params      : gin.Context
return      : context.Next
Author      : ssjpooh
Date        : 2023.10.10
*/
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
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid_token"})
			return
		}

	} else {

		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid_token"})
		return
	}

	// 여기 까지자 기본 토큰 validation 체크

	// 여기서 부터 추가 검증 ( count / 및 last request date )

	auth := memory.GetAuthInfo(bearerToken)

	if auth.CallCount > 50 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "too many call error"})
		return
	}

	log.Println("time.Now().Unix() : ", time.Now().Unix())
	log.Println("auth.LastRequestDt : ", auth.LastRequestDt)
	log.Println("calc : ", time.Now().Unix()-auth.LastRequestDt)
	if time.Now().Unix()-auth.LastRequestDt < 1000 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "too many call error"})
		return
	}

	// 호출 카운트 추가 및 마지막 호출 시간 초기화
	auth.CallCount = auth.CallCount + 1
	auth.LastRequestDt = time.Now().Unix()
	memory.SetAuthInfo(bearerToken, auth.ClientId, auth.ServerAddr, auth.CallCount, auth.ExpiredDt, auth.LastRequestDt)

	c.Next()
}
