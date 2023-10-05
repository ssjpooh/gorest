package auth

type ClientDetails struct {
	ClientID     string `db:"client_id"`
	ClientSecret string `db:"client_secret"`
}

type OauthInfo struct {
	ClientID   string `db:"client_id"`
	ExpiresAT  int64  `db:"expires_at"`
	Token      string `db:"token"`
	RFToken    string `db:"refresh_tokne"`
	ServerAddr string `db:"server_address"`
}

var JWTKey = []byte("FOXEduP@ssW0rd")
