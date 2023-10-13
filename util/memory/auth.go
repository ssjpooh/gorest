package memory

import (
	"fmt"
	oauthInfo "restApi/model/auth"
	dbHandler "restApi/util/db"
	"time"

	logger "restApi/util/log"
)

type AuthInfoMap map[string]oauthInfo.AuthInfo

var GlobalAuthInfoMap AuthInfoMap

/*
Description : auth 실행 할때 초기화
Params      :
return      :
Author      : ssjpooh
Date        : 2023.10.10
*/
func Init() {
	GlobalAuthInfoMap = make(AuthInfoMap)
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		for range ticker.C {
			checkGlobalMap()
		}
	}()
	logger.Logger(logger.GetFuncNm(), "Global Map setting Success")
}

/*
Description : 시간 마다 global map 의 expire 을 체크 하여 인증 시간이 만료된 토큰은 삭제
Params      :
return      :
Author      : ssjpooh
Date        : 2023.10.10
*/
func checkGlobalMap() {
	logger.Logger(logger.GetFuncNm(), "is Alive")
	for key := range GlobalAuthInfoMap {
		if GlobalAuthInfoMap[key].ExpiredDt < time.Now().Unix() {
			logger.Logger(logger.GetFuncNm(), fmt.Sprintf("expire date : %d , now : %d", GlobalAuthInfoMap[key].ExpiredDt, time.Now().Unix()))
			DelAuthInfo(key)
		}
	}
}

/*
Description : 나의 서버의 globalMap 에 key 값이 token 데이터를 찾아 없으면 db 에서 조회 하여 추가
Params      : token
return      : oauthInfo.AuthInfo
Author      : ssjpooh
Date        : 2023.10.10
*/
func GetAuthInfo(token string) oauthInfo.AuthInfo {

	if authInfo, exists := GlobalAuthInfoMap[token]; exists {
		return authInfo
	} else {
		logger.Logger(logger.GetFuncNm(), "GetAuthInfo not exitst")
		var oauth oauthInfo.OauthInfo
		err := dbHandler.Db.Get(&oauth, "SELECT refresh_token, client_id, expires_at, token, server_address from OAUTH_CLIENT_TOKENS where token = ? ", token)
		logger.Logger(logger.GetFuncNm(), "search token Info by token : ", token)
		if err != nil {
			logger.Logger(logger.GetFuncNm(), "select error :", err.Error())
		}

		SetAuthInfo(oauth.Token, oauth.ClientID, oauth.ServerAddr, 0, oauth.ExpiresAT, time.Now().Unix())

		return GlobalAuthInfoMap[token]
	}
}

/*
Description : global mpa 에 token set
Params      : token
Params      : clientId
Params      : serverAddr
Params      : callCount
Params      : expiredDt
Params      : lastRequestDt
return      :
Author      : ssjpooh
Date        : 2023.10.10
*/
func SetAuthInfo(token string, clientId string, serverAddr string, callCount int, expiredDt int64, lastRequestDt int64) {
	logger.Logger(logger.GetFuncNm(), "Set Global Map : ", token)
	GlobalAuthInfoMap[token] = oauthInfo.AuthInfo{ClientId: clientId, ServerAddr: serverAddr, CallCount: callCount, ExpiredDt: expiredDt, LastRequestDt: lastRequestDt}
}

/*
Description : global mpa 에 token 삭제
Params      : token
return      :
Author      : ssjpooh
Date        : 2023.10.10
*/
func DelAuthInfo(token string) {
	logger.Logger(logger.GetFuncNm(), "delete map key :", token)
	dbHandler.Db.Exec("DELETE FROM OAUTH_CLIENT_TOKENS WHERE TOKEN = ? ", token)
	delete(GlobalAuthInfoMap, token)
}

func PatchAuthInfo(beforeToken, newToken string) {
	logger.Logger(logger.GetFuncNm(), "patch before map key :", beforeToken, " new map key : ", newToken)
	if authInfo, exists := GlobalAuthInfoMap[beforeToken]; exists {
		GlobalAuthInfoMap[newToken] = oauthInfo.AuthInfo{ClientId: authInfo.ClientId, ServerAddr: authInfo.ServerAddr, CallCount: authInfo.CallCount, ExpiredDt: authInfo.ExpiredDt, LastRequestDt: authInfo.LastRequestDt}
		delete(GlobalAuthInfoMap, beforeToken)
	}
}
