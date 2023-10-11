package option

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/magiconair/properties"

	options "restApi/model/options"
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
			fmt.Println("Error getting executable path:", err)
			return
		}
		// 실행 파일의 디렉토리 경로 가져오기
		projectDir = filepath.Dir(exePath)
	} else {

		projectDir = filepath.Join(gopath, "src")
	}

	// .properties 파일 경로 설정
	configFile := filepath.Join(projectDir, "options.properties")

	log.Println("configFile : ", configFile)
	// .properties 파일 읽기
	p, err := properties.LoadFile(configFile, properties.UTF8)
	if err != nil {
		fmt.Println("Error reading .properties file:", err)
	}

	// 설정 값 읽기
	Prop.Url = p.GetString("db_url", "") // 기본값을 ""로 설정
	Prop.Id = p.GetString("db_id", "")
	Prop.Pw = p.GetString("db_pw", "")
	Prop.Name = p.GetString("db_nm", "")
}
