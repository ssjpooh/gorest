package auth

type ClientDetails struct {
	ClientID     string `db:"client_id"`
	ClientSecret string `db:"client_secret"`
}

type OauthInfo struct {
	ClientID  string `db:"client_id"`
	ExpiresAT int64  `db:"expires_at"`
	Token     string `db:"token"`
}
