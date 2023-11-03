package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"

	mobiletokenapi "restApi/controller/api/mobileTokenApi"
	tokenapi "restApi/controller/api/tokenApi"
	memoryApi "restApi/controller/authApi/memoryApi"
	roomsApi "restApi/controller/authApi/roomsApi"
	sitesApi "restApi/controller/authApi/sitesApi"
	usersApi "restApi/controller/authApi/usersApi"

	auth "restApi/model/auth"
	dbHandler "restApi/util/db"

	_ "restApi/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	gmap "restApi/util/memory"

	logger "restApi/util/log"
	options "restApi/util/options"

	vbase "restApi/model/vbase"

	"github.com/gin-contrib/cors"
)

var logFile = time.Now().Format("2006-01-02 15:04")

// @title Swagger Example API
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/
// @BasePath /v1

// @SecurityDefinitions.apikey Authorization
// @type apiKey
// @in header
// @name Authorization
func main() {

	logger.SetFileName(logFile)

	logger.Logger(logger.GetFuncNm(), "main Start")
	// map 설정
	gmap.Init()

	// option 설정
	options.Init()

	dbHandler.DbConnect()
	defer dbHandler.Db.Close()

	go func() {
		fs := http.FileServer(http.Dir("./web"))
		http.Handle("/", fs)
		// http.Handle("/login", http.HandlerFunc(loginHandler))
		http.Handle("/roomList", http.HandlerFunc(roomListHandler))
		http.ListenAndServeTLS(":443", options.Prop.CrtPath, options.Prop.KeyPath, nil)

	}()
	router := gin.Default()

	// CORS 설정
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Authorization", "Content-Type"}

	router.Use(cors.New(config))
	// v1 그룸에 포함 되지 않는 api들
	router.GET("/", indexHandler)
	router.POST("/login", apiLoginHandler)

	v1 := router.Group("/v1")
	// 토큰
	tokenapi.TokenApiHandler(v1)
	// 모바일 토큰
	mobiletokenapi.MobileTokenApiHandler(v1)

	// users
	sitesApi.SitesApiHandler(v1)
	// sub users
	usersApi.UsersApiHandler(v1)
	// rooms
	roomsApi.RoomApiHandler(v1)

	// 인증 토큰  관리
	memoryApi.MemoryApiHaneler(v1)

	url := ginSwagger.URL("https://local.foxedu.kr:443/swagger/doc.json") // The url pointing to API definition
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	err := router.RunTLS(":1443", options.Prop.CrtPath, options.Prop.KeyPath)

	if err != nil {
		panic(err)
	}

}

// 일반 웹

func roomListHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/index.html")
}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	getId := r.PostFormValue("id")
	getPassword := r.PostFormValue("password")

	logger.Logger(logger.GetFuncNm(), " getId : ", getId, " getPassword : ", getPassword)

	data := url.Values{}
	data.Set("id", getId)
	data.Set("password", getPassword)
	w.Header().Set("Content-Type", "application/json")
	response, err := http.PostForm("https://local.foxedu.kr:1443/login", data)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return
	}
	defer response.Body.Close()

	logger.Logger(logger.GetFuncNm(), " response : ", response.Status)
	// 상태 코드 확인
	if response.StatusCode != http.StatusOK {
		fmt.Printf("Server returned status code %d: %s\n", response.StatusCode, response.Status)
		return
	}

	// 응답 본문(body) 읽기
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %s\n", err)
		return
	}

	// 응답 본문을 문자열로 변환하여 출력
	bodyString := string(bodyBytes)
	w.Write([]byte(bodyString))
	// 필요한 경우, JSON 응답을 파싱할 수도 있습니다.
	// 예: json.Unmarshal(bodyBytes, &yourStruct)
}

// description : 접속 확인
func indexHandler(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, "ok")
}

// api login
func apiLoginHandler(context *gin.Context) {
	id := context.PostForm("id")
	pw := context.PostForm("password")
	logger.Logger(logger.GetFuncNm(), "id : ", id, " pw : ", pw)

	var userInfo vbase.Users
	userInfo = usersApi.GetUsersInfoById(context, id)

	if userInfo.UserID == "" {
		context.JSON(http.StatusNotFound, gin.H{"message": "please check id or password"})
	}
	var matchPassword bool
	if options.Prop.DebugMode {
		if pw == userInfo.Passwd {
			matchPassword = true
		} else {
			matchPassword = false
		}
	} else {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)

		if pw == string(hashedPassword) {
			matchPassword = true
		} else {
			matchPassword = false
		}
	}

	if !matchPassword {
		context.JSON(http.StatusNotFound, gin.H{"message": "please check id or password"})
	}
	var siteIdx string
	siteIdx = userInfo.SiteIdx

	var siteInfo vbase.Sites
	siteInfo = sitesApi.GetSitesInfoByIdx(context, siteIdx)

	if siteInfo.SiteID == "" {
		context.JSON(http.StatusNotFound, gin.H{"message": "please check id or password"})
	}

	var oauthClientDetails auth.OAuthClientDetails
	oauthClientDetails = tokenapi.GetOauthClientDetails(context, siteIdx)
	if oauthClientDetails.ClientID == "" {
		context.JSON(http.StatusNotFound, gin.H{"message": "not registed client info"})
	}

	clientId := oauthClientDetails.ClientID
	clientSecret := oauthClientDetails.ClientSecret

	tokenapi.TokenGlobalApi(context, clientId, clientSecret)

}
