package auth

import "database/sql"

var JWTKey = []byte("FOXEduP@ssW0rd")

type AuthInfo struct {
	ClientId      string
	ExpiredDt     int64
	LastRequestDt int64
	CallCount     int
	ServerAddr    sql.NullString
	ApiName       string
}
