package util

import (
	"database/sql"
	"net"

	logger "restApi/util/log"

	"github.com/google/uuid"
)

func GenterateUUID() string {

	uuidObj, err := uuid.NewRandom()
	if err != nil {
		logger.Logger(logger.GetFuncNm(), "UUID generate error : ", err.Error())
	}

	return uuidObj.String()
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func GetStringOrNil(s sql.NullString) *string {
	if !s.Valid {
		// SQL NULL을 나타냄
		return nil
	}
	// 문자열에 대한 포인터를 반환함
	return &s.String
}
