package memory

import (
	"log"
	outhInfo "restApi/model/auth"
	"time"
)

type AuthInfoMap map[string]outhInfo.AuthInfo

var GlobalAuthInfoMap AuthInfoMap

func Init() {

	GlobalAuthInfoMap = make(AuthInfoMap)
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		for range ticker.C {
			checkGlobalMap(GlobalAuthInfoMap)
		}
	}()

}

func checkGlobalMap(GlobalAuthInfoMap AuthInfoMap) {
	log.Println("AuthInfoMap : ", GlobalAuthInfoMap)
}

func GetAuthInfo(token string) outhInfo.AuthInfo {

	log.Println("get token : ", token)
	if authInfo, exists := GlobalAuthInfoMap[token]; exists {
		log.Println("exist")
		return authInfo
	} else {
		log.Println("not exist")
		return authInfo
	}
}

func SetAuthInfo(token string, clientId string, serverAddr string, callCount int, expiredDt int64, lastRequestDt int64) {
	log.Println("set token : ", token)
	GlobalAuthInfoMap[token] = outhInfo.AuthInfo{ClientId: clientId, ServerAddr: serverAddr, CallCount: callCount, ExpiredDt: expiredDt, LastRequestDt: lastRequestDt}
}

func DelAuthInfo(token string) {

}
