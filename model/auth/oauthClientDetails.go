package auth

import (
	"database/sql"
	"fmt"
	db "restApi/util/db"
	"strings"

	"gorm.io/gorm"
)

type OAuthClientDetails struct {
	gorm.Model
	SiteIdx              string         `gorm:"type:char(36);not null;index" db:"site_idx"`
	ClientID             string         `gorm:"type:char(36);not null;primaryKey" db:"client_id"`
	ClientSecret         string         `gorm:"type:char(36);not null" db:"client_secret"`
	Scope                sql.NullString `gorm:"type:varchar(255)" db:"scope"`                   // Nullable
	AuthorizedGrantTypes sql.NullString `gorm:"type:varchar(255)" db:"authorized_grant_types"`  // Nullable
	WebServerRedirectURI sql.NullString `gorm:"type:varchar(255)" db:"web_server_redirect_uri"` // Nullable
	CDate                sql.NullString `gorm:"type:varchar(14);index" db:"cdate"`              // Nullable
	MDate                sql.NullString `gorm:"type:varchar(14);index" db:"mdate"`              // Nullable
}

func (OAuthClientDetails) TableName() string {
	return "oauth_client_details"
}

var OAuthClientDetailsColumns = fmt.Sprintf(" %s ", strings.Join(db.ColumnsForStruct(OAuthClientDetails{}), ", "))
