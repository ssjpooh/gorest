package util

import (
	"fmt"
	"log"
	"net"

	"github.com/google/uuid"
)

func GenterateUUID() uuid.UUID {

	uuidObj, err := uuid.NewRandom()
	if err != nil {
		fmt.Println("UUID 생성 중 오류 발생:", err)
		log.Print(err)
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
