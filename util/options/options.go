package option

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/magiconair/properties"

	options "restApi/model/options"
	logger "restApi/util/log"
)

var Prop options.OptionsInfo

func Init() {

	var projectDir string
	// 실행 파일의 경로 가져오기

	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		// 프로젝트 루트 디렉토리 경로 설정
		exePath, err := os.Executable()
		if err != nil {
			logger.Logger(logger.GetFuncNm(), "[err],  Error getting executable path:", err.Error())
			return
		}
		// 실행 파일의 디렉토리 경로 가져오기
		projectDir = filepath.Dir(exePath)
	} else {

		projectDir = filepath.Join(gopath, "src")
	}

	// .properties 파일 경로 설정
	configFile := filepath.Join(projectDir, "options.properties")

	logger.Logger(logger.GetFuncNm(), "optionFilePath : ", configFile)
	// .properties 파일 읽기
	p, err := properties.LoadFile(configFile, properties.UTF8)
	if err != nil {
		logger.Logger(logger.GetFuncNm(), "[err] Error reading .properties file:", err.Error())
	}

	// 설정 값 읽기
	Prop.Url = strings.TrimSpace(p.GetString("db_url", "")) // 기본값을 ""로 설정
	Prop.Id = strings.TrimSpace(p.GetString("db_id", ""))
	Prop.Pw = strings.TrimSpace(p.GetString("db_pw", ""))
	Prop.Name = strings.TrimSpace(p.GetString("db_nm", ""))
	Prop.CrtPath = strings.TrimSpace(p.GetString("ssl_crt", ""))
	Prop.KeyPath = strings.TrimSpace(p.GetString("ssl_key", ""))

	logger.Logger(logger.GetFuncNm(), "option Init Success")
}
