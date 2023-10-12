package util

import (
	"net"

	logger "restApi/util/log"

	"github.com/google/uuid"
)

func GenterateUUID() uuid.UUID {

	uuidObj, err := uuid.NewRandom()
	if err != nil {
		logger.Logger(logger.GetFuncNm(), "UUID generate error : ", err.Error())
	}

	return uuidObj
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
