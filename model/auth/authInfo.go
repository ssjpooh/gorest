package auth

var JWTKey = []byte("FOXEduP@ssW0rd")

type AuthInfo struct {
	ClientId      string
	ExpiredDt     int64
	LastRequestDt int64
	CallCount     int
	ServerAddr    string
	ApiName       string
}
