package util

import (
	"fmt"
	"log"

	"github.com/google/uuid"
)

func GenterateUUID() uuid.UUID {

	uuidObj, err := uuid.NewRandom()
	if err != nil {
		fmt.Println("UUID 생성 중 오류 발생:", err)
		log.Fatal(err)
	}

	return uuidObj
}
