package subusersapi

import (
	"database/sql"
	"net/http"
	users "restApi/model/vbase"

	dbHandler "restApi/util/db"
	logger "restApi/util/log"

	authHandler "restApi/util/auth"

	utils "restApi/util"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func getUserList(context *gin.Context) []users.Users {

	var userList []users.Users

	query := dbHandler.MakeQuery(dbHandler.SELECT, users.UsersColumns, dbHandler.FROM, "users")

	err := dbHandler.Db.Get(&userList, query)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "not Found"})
		logger.Logger(logger.GetFuncNm(), "select error ", err.Error())
	}

	return userList
}

func GetUsersInfoById(context *gin.Context, id string) users.Users {
	var subUserInfo users.Users
	query := dbHandler.MakeQuery(dbHandler.SELECT, users.UsersColumns, dbHandler.FROM, "users", dbHandler.WHERE, "user_id = ? ")
	logger.Logger(logger.GetFuncNm(), " query : ", query)
	err := dbHandler.Db.Get(&subUserInfo, query, id)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "not Found"})
		logger.Logger(logger.GetFuncNm(), "select error :", err.Error())
		return subUserInfo
	}

	return subUserInfo

}

func inserUsersInfo(context *gin.Context) sql.Result {

	var newUser users.Users
	err := context.ShouldBindJSON(&newUser)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		logger.Logger(logger.GetFuncNm(), "insert error : ", err.Error())
	}

	// is exist sub users info
	subUserInfo := GetUsersInfoById(context, newUser.UserID)
	if subUserInfo.UserID != "" {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		logger.Logger(logger.GetFuncNm(), "insert error : ", err.Error())
	}

	query := "INSERT INTO USER_TBL (owner_idx, user_id, user_passwd, kor_user_name, eng_user_name ) values (?, ? , ? , ? , ?)"

	// 비밀번호를 해싱
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Passwd), bcrypt.DefaultCost)
	if err != nil {
		logger.Logger(logger.GetFuncNm(), " password hash error : ", err.Error())
	}

	// 무작위 UUID 생성
	ownerIdx := utils.GenterateUUID()
	result, err := dbHandler.Db.Exec(query, ownerIdx, newUser.UserID, hashedPassword, newUser.Name)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		logger.Logger(logger.GetFuncNm(), "insert error : ", err.Error())
	}

	return result

}

func patchUsersInfo(context *gin.Context, id string) users.Users {

	subUserInfo := GetUsersInfoById(context, id)

	var newSubUser users.Users
	err := context.ShouldBindJSON(&newSubUser)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		logger.Logger(logger.GetFuncNm(), " patch error : ", err.Error())
	}

	if subUserInfo.UserID == "" {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad param"})
	} else {
		if subUserInfo.Name != "" {
			subUserInfo.Name = subUserInfo.Name
		}

		query := "UPDATE USER_TBL set kor_user_name = ? , eng_user_name = ? where user_id = ? "
		_, err := dbHandler.Db.Exec(query, subUserInfo.Name, subUserInfo.UserID)
		if err != nil {
			context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
			logger.Logger(logger.GetFuncNm(), "patch error : ", err.Error())
		}

	}
	return subUserInfo

}

func deleteUsersInfo(context *gin.Context, id string) sql.Result {

	subUserInfo := GetUsersInfoById(context, id)

	var newSubUser users.Users
	err := context.ShouldBindJSON(&newSubUser)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		logger.Logger(logger.GetFuncNm(), " patch error : ", err.Error())
	}

	if subUserInfo.UserID == "" {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad param"})
	}

	query := "DELETE FROM  USER_TBL set kor_user_name = ? , eng_user_name = ? where user_id = ? "
	result, err := dbHandler.Db.Exec(query, subUserInfo.Name, subUserInfo.UserID)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		logger.Logger(logger.GetFuncNm(), "patch error : ", err.Error())
	}

	return result
}

func UsersApiHandler(v1 *gin.RouterGroup) {

	v1.GET("/subUsers", authHandler.Authenticate, func(context *gin.Context) {
		UserList := getUserList(context)
		context.IndentedJSON(http.StatusOK, UserList)
	})

	// @Summary Show an account
	// @Description Get account by ID
	// @Tags memgers
	// @Accept  json
	// @Produce  json
	// @Param member body Memeber true "Member"
	// @Success 200 {object} users.Users
	// @Router /members/:id [post]
	v1.GET("/subUsers/:userSID", authHandler.Authenticate, func(context *gin.Context) {
		id := context.Param("userSID")
		logger.Logger(logger.GetFuncNm(), "id : ", id)
		result := GetUsersInfoById(context, id)
		context.IndentedJSON(http.StatusOK, result)
	})

	// @Summary Show an account
	// @Description Get account by ID
	// @Tags memgers
	// @Accept  json
	// @Produce  json
	// @Param member body Memeber true "Member"
	// @Success 200 {object} users.Users
	// @Router /members/:id [post]
	v1.POST("/subUsers", authHandler.Authenticate, func(context *gin.Context) {
		result := inserUsersInfo(context)
		context.IndentedJSON(http.StatusCreated, result)
	})

	// @Summary Show an account
	// @Description Get account by ID
	// @Tags memgers
	// @Accept  json
	// @Produce  json
	// @Param	id path string true "ID"
	// @Success 200 {object} users.Users
	// @Router /members/:id [patch]
	v1.PATCH("/subUsers/:id", authHandler.Authenticate, func(context *gin.Context) {
		id := context.Param("userId")
		userInfo := patchUsersInfo(context, id)
		context.IndentedJSON(http.StatusOK, userInfo)
	})

	// @Summary Show an account
	// @Description Get account by ID
	// @Tags memgers
	// @Accept  json
	// @Produce  json
	// @Param	id path string true "ID"
	// @Success 200 {object} users.Users
	// @Router /members/:id [patch]
	v1.DELETE("/subUsers/:id", authHandler.Authenticate, func(context *gin.Context) {
		id := context.Param("userId")
		userInfo := deleteUsersInfo(context, id)
		context.IndentedJSON(http.StatusOK, userInfo)
	})

}
