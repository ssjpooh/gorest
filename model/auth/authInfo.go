package auth

import (
	"fmt"
	db "restApi/util/db"
	"strings"
)

type ClientDetails struct {
	ClientID     string `db:"CLIENT_ID"`
	OwnerIdx     string `db:"OWNER_IDX"`
	ClientSecret string `db:"CLIENT_SECRET"`
}

type OauthInfo struct {
	ClientID   string `db:"client_id"`
	ExpiresAT  int64  `db:"expires_at"`
	Token      string `db:"token"`
	RFToken    string `db:"refresh_token"`
	ServerAddr string `db:"server_address"`
}

var JWTKey = []byte("FOXEduP@ssW0rd")

type AuthInfo struct {
	ClientId      string
	ExpiredDt     int64
	LastRequestDt int64
	CallCount     int
	ServerAddr    string
	ApiName       string
}

var OAuthClientDetailsColumns = fmt.Sprintf(" %s ", strings.Join(db.ColumnsForStruct(ClientDetails{}), ", "))
var OAuthClientTokensColumns = fmt.Sprintf(" %s ", strings.Join(db.ColumnsForStruct(OauthInfo{}), ", "))
