package auth

import (
	"fmt"
	db "restApi/util/db"
	"strings"

	"gorm.io/gorm"
)

type OAuthClientTokens struct {
	gorm.Model
	ClientID      string `gorm:"type:char(36);not null;primaryKey;index;foreignKey:OAuthClientDetailsRefer" db:"client_id"`
	ExpiresAt     int64  `gorm:"not null" db:"expires_at"`
	Token         string `gorm:"type:varchar(1000);not null" db:"token"`
	RefreshToken  string `gorm:"type:varchar(1000);not null" db:"refresh_token"`
	ServerAddress string `gorm:"type:varchar(50)" db:"server_address"` // Nullable
	CDate         string `gorm:"type:varchar(14);index" db:"cdate"`    // Nullable
	MDate         string `gorm:"type:varchar(14);index" db:"mdate"`    // Nullable
}

func (OAuthClientTokens) TableName() string {
	return "oauth_client_tokens"
}

var OAuthClientTokensColumns = fmt.Sprintf(" %s ", strings.Join(db.ColumnsForStruct(OAuthClientTokens{}), ", "))
