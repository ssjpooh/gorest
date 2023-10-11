package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func Logger(method string, infos ...string) {

	var logDir string
	// 실행 파일의 경로 가져오기

	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		// 프로젝트 루트 디렉토리 경로 설정
		exePath, err := os.Executable()
		if err != nil {
			fmt.Println("Error getting executable path:", err)
			return
		}
		// 실행 파일의 디렉토리 경로 가져오기
		logDir = filepath.Dir(exePath + "/logs")
	} else {

		logDir = filepath.Join(gopath, "src", "logs")
	}

	// 폴더가 존재하지 않으면 생성
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.Mkdir(logDir, 0755)
	}

	// 현재 날짜를 YYYY-MM-DD 형식으로 가져옵니다.
	currentDate := time.Now().Format("2006-01-02 15:04")

	// 로그 파일 이름을 현재 날짜로 설정합니다. 경로를 포함합니다.
	logFilePath := filepath.Join(logDir, fmt.Sprintf("[restApi]_%s.log", currentDate))

	// 파일을 생성 또는 열기
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 로그 출력을 해당 파일로 설정
	log.SetOutput(file)

	// 로그 메시지 작성
	log.Printf("[" + method + "] " + strings.Join(infos, "") + "\n")

}

func GetFuncNm() string {
	pc, _, _, _ := runtime.Caller(1) // 1은 현재 함수의 호출자를 나타냅니다.
	return runtime.FuncForPC(pc).Name()
}
