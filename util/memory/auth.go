package memory

import (
	"log"
	oauthInfo "restApi/model/auth"
	dbHandler "restApi/util/db"
	"time"
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
}

/*
Description : 시간 마다 global map 의 expire 을 체크 하여 인증 시간이 만료된 토큰은 삭제
Params      :
return      :
Author      : ssjpooh
Date        : 2023.10.10
*/
func checkGlobalMap() {
	for key := range GlobalAuthInfoMap {
		if GlobalAuthInfoMap[key].ExpiredDt < time.Now().Unix() {
			delete(GlobalAuthInfoMap, key)
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

	log.Println("get token : ", token)
	if authInfo, exists := GlobalAuthInfoMap[token]; exists {
		log.Println("::::::::: exist :::::::::")
		return authInfo
	} else {
		log.Println("::::::::: not exist :::::::::")
		var oauth oauthInfo.OauthInfo
		err := dbHandler.Db.Get(&oauth, "SELECT refresh_token, client_id, expires_at, token, server_address from oauth_tokens where token = ? ", token)

		if err != nil {
			log.Println("4444444444444444")
			SetAuthInfo(oauth.Token, oauth.ClientID, oauth.ServerAddr, 0, oauth.ExpiresAT, time.Now().Unix())
		}

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
	log.Println("222222222222222")
	log.Println("set token : ", token)
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

	delete(GlobalAuthInfoMap, token)
}

func PatchAuthInfo(beforeToken, newToken string) {

	log.Println("get token : ", beforeToken)
	if authInfo, exists := GlobalAuthInfoMap[beforeToken]; exists {
		log.Println("::::::::: exist :::::::::")
		GlobalAuthInfoMap[newToken] = oauthInfo.AuthInfo{ClientId: authInfo.ClientId, ServerAddr: authInfo.ServerAddr, CallCount: authInfo.CallCount, ExpiredDt: authInfo.ExpiredDt, LastRequestDt: authInfo.LastRequestDt}
		delete(GlobalAuthInfoMap, beforeToken)
	}
}
