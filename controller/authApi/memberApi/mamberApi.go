package member

import (
	"database/sql"
	"net/http"

	"log"
	members "restApi/model/members"

	dbHandler "restApi/util/db"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"

	utils "restApi/util"

	authHandler "restApi/util/auth"
)

func getMemberList(context *gin.Context) []members.Member {

	var userList []members.Member

	query := "SELECT owner_idx, user_id, user_passwd, kor_user_name, eng_user_name FROM user_tbl"

	err := dbHandler.Db.Select(&userList, query)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "not Found"})
		log.Print(err)
	}

	return userList

}

func getMemberInfo(context *gin.Context, id string) members.Member {
	var userInfo members.Member
	query := "SELECT owner_idx, user_id, user_passwd, kor_user_name, eng_user_name FROM user_tbl WHERE user_id = ?"

	err := dbHandler.Db.Get(&userInfo, query, id)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "not Found"})
		log.Print(err)
	}

	return userInfo

}
func inserMemberInfo(context *gin.Context) sql.Result {

	var newUser members.Member
	err := context.ShouldBindJSON(&newUser)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		log.Print(err)
	}

	query := "INSERT INTO USER_TBL (owner_idx, user_id, user_passwd, kor_user_name, eng_user_name ) values (?, ? , ? , ? , ?)"

	// 비밀번호를 해싱
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.PASSWD), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	// 무작위 UUID 생성
	ownerIdx := utils.GenterateUUID()
	result, err := dbHandler.Db.Exec(query, ownerIdx, newUser.ID, hashedPassword, newUser.KORName, newUser.ENGName)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		log.Print(err)
	}

	insertOauth(context, ownerIdx)

	return result

}

func insertOauth(context *gin.Context, ownerIdx uuid.UUID) sql.Result {

	clientId := utils.GenterateUUID()
	clientSecret := utils.GenterateUUID()

	query2 := "INSERT INTO OAUTH_CLIENT_DETAILS ( owner_idx, client_id, client_secret) values ( ?, ?, ? )"

	result, err := dbHandler.Db.Exec(query2, ownerIdx, clientId, clientSecret)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		log.Print(err)
	}

	return result
}

func MmberApiHandler(v1 *gin.RouterGroup) {

	v1.GET("/members", authHandler.Authenticate, func(context *gin.Context) {
		userList := getMemberList(context)
		context.IndentedJSON(http.StatusOK, userList)
	})

	// v1.GET("/members", func(context *gin.Context) {
	// 	fmt.Println("TEST 3")
	// 	userList := getMemberList(context)
	// 	context.IndentedJSON(http.StatusOK, userList)
	// })

	// v1.GET("/members/:id", authHandler.Authenticate, func(context *gin.Context) {
	// 	id := context.Param("id")
	// 	userInfo := getMemberInfo(context, id)
	// 	context.IndentedJSON(http.StatusOK, userInfo)
	// })

	v1.GET("/members/:id", func(context *gin.Context) {
		id := context.Param("id")
		userInfo := getMemberInfo(context, id)
		context.IndentedJSON(http.StatusOK, userInfo)
	})
	v1.POST("/members", func(context *gin.Context) {
		result := inserMemberInfo(context)
		context.IndentedJSON(http.StatusCreated, result)
	})

}
