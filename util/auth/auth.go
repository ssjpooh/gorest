package auth

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	outhInfo "restApi/model/auth"
	memory "restApi/util/memory"

	"github.com/dgrijalva/jwt-go"

	logger "restApi/util/log"
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

	auth, err := memory.GetAuthInfo(bearerToken, c.Request.RequestURI)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid_token"})
		return
	}

	if auth.CallCount > 50 {
		logger.Logger(logger.GetFuncNm(), strconv.Itoa(auth.CallCount))
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "too many call by call Count error"})
		return
	}

	logger.Logger(logger.GetFuncNm(), " data :  request uri : ", c.Request.RequestURI, "  auth.ApiName : ", auth.ApiName)
	if time.Now().Unix()-auth.LastRequestDt <= 2 && auth.ApiName == c.Request.RequestURI {
		logger.Logger(logger.GetFuncNm(), strconv.FormatInt(time.Now().Unix(), 10))
		logger.Logger(logger.GetFuncNm(), strconv.FormatInt(auth.LastRequestDt, 10))
		logger.Logger(logger.GetFuncNm(), strconv.FormatInt(time.Now().Unix()-auth.LastRequestDt, 10))
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "too earily call by lastRequestd Date error"})
		return
	}

	// 호출 카운트 추가 및 마지막 호출 시간 초기화
	auth.CallCount = auth.CallCount + 1
	auth.LastRequestDt = time.Now().Unix()
	auth.ApiName = c.Request.RequestURI
	memory.SetAuthInfo(bearerToken, auth.ClientId, auth.ServerAddr, auth.CallCount, auth.ExpiredDt, auth.LastRequestDt, auth.ApiName)

	c.Next()
}
