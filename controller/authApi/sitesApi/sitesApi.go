package usersapi

import (
	"database/sql"
	"net/http"

	sites "restApi/model/vbase"

	dbHandler "restApi/util/db"
	logger "restApi/util/log"

	"github.com/gin-gonic/gin"

	"golang.org/x/crypto/bcrypt"

	utils "restApi/util"

	authHandler "restApi/util/auth"
)

/*
Description : get member list
Params      : gin.Context
return      : []memvers.Member
Author      : ssjpooh
Date        : 2023.10.10
*/
func getSitesList(context *gin.Context) []sites.Sites {

	var userList []sites.Sites

	query := dbHandler.MakeQuery(dbHandler.SELECT, sites.SitesColumns, dbHandler.FROM, "sites")

	err := dbHandler.Db.Select(&userList, query)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "not Found"})
		logger.Logger(logger.GetFuncNm(), "select error ", err.Error())
	}

	return userList

}

func GetSitesInfoById(context *gin.Context, id string) sites.Sites {
	var userInfo sites.Sites
	query := dbHandler.MakeQuery(dbHandler.SELECT, sites.SitesColumns, dbHandler.FROM, "sites", dbHandler.WHERE, "site_id = ? ")

	err := dbHandler.Db.Get(&userInfo, query, id)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "not Found"})
		logger.Logger(logger.GetFuncNm(), "select error :", err.Error())
		return userInfo
	}

	return userInfo

}

func GetSitesInfoByIdx(context *gin.Context, idx string) sites.Sites {
	var userInfo sites.Sites
	query := dbHandler.MakeQuery(dbHandler.SELECT, sites.SitesColumns, dbHandler.FROM, "sites", dbHandler.WHERE, "site_idx = ? ")
	logger.Logger(logger.GetFuncNm(), " query : ", query, " idx : ", idx)
	err := dbHandler.Db.Get(&userInfo, query, idx)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "not Found"})
		logger.Logger(logger.GetFuncNm(), "select error :", err.Error())
		return userInfo
	}

	return userInfo

}

func inserSitesInfo(context *gin.Context) sql.Result {

	var newSites sites.Sites
	err := context.ShouldBindJSON(&newSites)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		logger.Logger(logger.GetFuncNm(), "insert error : ", err.Error())
	}

	siteInfo := GetSitesInfoById(context, newSites.SiteID)

	// is exist users info
	if siteInfo.SiteID != "" {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		logger.Logger(logger.GetFuncNm(), "insert error : ", err.Error())
	}

	query := "INSERT INTO USER_TBL (owner_idx, user_id, user_passwd, kor_user_name, eng_user_name ) values (?, ? , ? , ? , ?)"

	var hashedPassword []byte
	// 비밀번호를 해싱
	if newSites.Passwd.Valid {

		hashedPassword, err = bcrypt.GenerateFromPassword([]byte(newSites.Passwd.String), bcrypt.DefaultCost)
		if err != nil {
			logger.Logger(logger.GetFuncNm(), " password hash error : ", err.Error())
		}
	} else {
		hashedPassword = nil

	}
	// 무작위 UUID 생성
	ownerIdx := utils.GenterateUUID()
	result, err := dbHandler.Db.Exec(query, ownerIdx, newSites.SiteID, hashedPassword, newSites.Name)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		logger.Logger(logger.GetFuncNm(), "insert error : ", err.Error())
	}

	insertOauth(context, ownerIdx)

	return result

}

func insertOauth(context *gin.Context, ownerIdx string) sql.Result {

	clientId := utils.GenterateUUID()
	clientSecret := utils.GenterateUUID()

	query := "INSERT INTO OAUTH_CLIENT_DETAILS ( owner_idx, client_id, client_secret) values ( ?, ?, ? )"

	result, err := dbHandler.Db.Exec(query, ownerIdx, clientId, clientSecret)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		logger.Logger(logger.GetFuncNm(), "insert error : ", err.Error())
	}

	return result
}

func patchSitesInfo(context *gin.Context, siteId string) sites.Sites {

	siteInfo := GetSitesInfoById(context, siteId)

	var newSites sites.Sites
	err := context.ShouldBindJSON(&newSites)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		logger.Logger(logger.GetFuncNm(), " patch error : ", err.Error())
	}

	if siteInfo.SiteID == "" {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad param"})
	} else {
		if newSites.Name != "" {
			siteInfo.Name = newSites.Name
		}

		query := "UPDATE USER_TBL set kor_user_name = ? , eng_user_name = ? where site_id = ? "
		_, err := dbHandler.Db.Exec(query, siteInfo.Name, siteInfo.SiteID)
		if err != nil {
			context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
			logger.Logger(logger.GetFuncNm(), "patch error : ", err.Error())
		}

	}
	return siteInfo

}

func deleteSitesInfo(context *gin.Context, siteId string) sql.Result {

	siteInfo := GetSitesInfoById(context, siteId)

	var newSites sites.Sites
	err := context.ShouldBindJSON(&newSites)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		logger.Logger(logger.GetFuncNm(), " patch error : ", err.Error())
	}

	if siteInfo.SiteID == "" {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad param"})
	}

	query := "UPDATE USER_TBL set kor_user_name = ? , eng_user_name = ? where site_id = ? "
	result, err := dbHandler.Db.Exec(query, siteInfo.Name, siteInfo.SiteID)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		logger.Logger(logger.GetFuncNm(), "patch error : ", err.Error())
	}

	return result

}

func SitesApiHandler(v1 *gin.RouterGroup) {

	v1.GET("/users", authHandler.Authenticate, func(context *gin.Context) {
		userList := getSitesList(context)
		context.IndentedJSON(http.StatusOK, userList)
	})

	// @Summary Show an account
	// @Description Get member by id
	// @Tags memgers
	// @Accept  json
	// @Produce  json
	// @Param	id path string true "ID"
	// @Success 200 {object} sites.Sites
	// @Router /members/:id [get]
	v1.GET("/users/:userID", authHandler.Authenticate, func(context *gin.Context) {
		id := context.Param("userID")
		userInfo := GetSitesInfoById(context, id)
		context.IndentedJSON(http.StatusOK, userInfo)
	})

	// @Summary Show an account
	// @Description Get account by ID
	// @Tags memgers
	// @Accept  json
	// @Produce  json
	// @Param	id path string true "ID"
	// @Success 200 {object} sites.Sites
	// @Router /members/:id [patch]
	v1.PATCH("/users/:id", authHandler.Authenticate, func(context *gin.Context) {
		id := context.Param("userId")
		userInfo := patchSitesInfo(context, id)
		context.IndentedJSON(http.StatusOK, userInfo)
	})

	// @Summary Show an account
	// @Description Get account by ID
	// @Tags memgers
	// @Accept  json
	// @Produce  json
	// @Param member body Memeber true "Member"
	// @Success 200 {object} sites.Sites
	// @Router /members/:id [post]
	v1.POST("/users", authHandler.Authenticate, func(context *gin.Context) {
		result := inserSitesInfo(context)
		context.IndentedJSON(http.StatusCreated, result)
	})

	// @Summary Show an account
	// @Description Get account by ID
	// @Tags memgers
	// @Accept  json
	// @Produce  json
	// @Param member body Memeber true "Member"
	// @Success 200 {object} sites.Sites
	// @Router /members/:id [post]
	v1.DELETE("/users/:id", authHandler.Authenticate, func(context *gin.Context) {
		id := context.Param("userId")
		result := deleteSitesInfo(context, id)
		context.IndentedJSON(http.StatusCreated, result)
	})

}
