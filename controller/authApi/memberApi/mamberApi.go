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

// @Summary Show an account
// @Description Get account list
// @Tags members
// @Accept  json
// @Produce  json
// @Success 200 {object} members.Member
// @Security Authorization
// @Router /members [get]
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
		return userInfo
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
		log.Print(err)
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

	query := "INSERT INTO OAUTH_CLIENT_DETAILS ( owner_idx, client_id, client_secret) values ( ?, ?, ? )"

	result, err := dbHandler.Db.Exec(query, ownerIdx, clientId, clientSecret)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		log.Print(err)
	}

	return result
}

func patchMemeberInfo(context *gin.Context, id string) members.Member {

	userInfo := getMemberInfo(context, id)

	var newUser members.Member
	err := context.ShouldBindJSON(&newUser)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		log.Print(err)
	}

	if userInfo.ID == "" {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad param"})
	} else {
		if newUser.KORName != "" {
			userInfo.KORName = newUser.KORName
		}

		if newUser.ENGName != "" {
			userInfo.ENGName = newUser.ENGName
		}

		query := "UPDATE USER_TBL set kor_user_name = ? , eng_user_name = ? where user_id = ? "
		_, err := dbHandler.Db.Exec(query, userInfo.KORName, userInfo.ENGName, userInfo.ID)
		if err != nil {
			context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
			log.Print(err)
		}

	}
	return userInfo

}

func MmberApiHandler(v1 *gin.RouterGroup) {

	v1.GET("/members", authHandler.Authenticate, func(context *gin.Context) {
		userList := getMemberList(context)
		context.IndentedJSON(http.StatusOK, userList)
	})

	// @Summary Show an account
	// @Description Get member by id
	// @Tags memgers
	// @Accept  json
	// @Produce  json
	// @Param	id path string true "ID"
	// @Success 200 {object} members.Member
	// @Router /members/:id [get]
	v1.GET("/members/:id", authHandler.Authenticate, func(context *gin.Context) {
		id := context.Param("id")
		userInfo := getMemberInfo(context, id)
		context.IndentedJSON(http.StatusOK, userInfo)
	})

	// @Summary Show an account
	// @Description Get account by ID
	// @Tags memgers
	// @Accept  json
	// @Produce  json
	// @Param	id path string true "ID"
	// @Success 200 {object} members.Member
	// @Router /members/:id [patch]
	v1.PATCH("/members/:id", authHandler.Authenticate, func(context *gin.Context) {
		id := context.Param("id")
		userInfo := patchMemeberInfo(context, id)
		context.IndentedJSON(http.StatusOK, userInfo)
	})

	// @Summary Show an account
	// @Description Get account by ID
	// @Tags memgers
	// @Accept  json
	// @Produce  json
	// @Param member body Memeber true "Member"
	// @Success 200 {object} members.Member
	// @Router /members/:id [post]
	v1.POST("/members", func(context *gin.Context) {
		result := inserMemberInfo(context)
		context.IndentedJSON(http.StatusCreated, result)
	})

}
