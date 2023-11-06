package vbase

import (
	"database/sql"
	"fmt"
	"restApi/util/db"
	"strings"
)

type Users struct {
	SiteIdx        string         `gorm:"type:char(36);not null;index" db:"site_idx"`
	UserIdx        string         `gorm:"type:char(36);primaryKey" db:"user_idx"`
	UserID         string         `gorm:"type:varchar(128);not null;unique" db:"user_id"`
	Passwd         sql.NullString `gorm:"type:varchar(32);not null" db:"passwd"`
	Email          string         `gorm:"type:varchar(128);default:'';not null;index" db:"email"`
	Name           string         `gorm:"type:nvarchar(64);not null;index" db:"name"`
	LoginServerIdx sql.NullString `gorm:"type:char(36);index" db:"login_server_idx"`
	LastLoginDate  sql.NullString `gorm:"type:varchar(14)" db:"last_login_date"`
	LastLogoutDate sql.NullString `gorm:"type:varchar(14)" db:"last_logout_date"`
	LastIPAddr     sql.NullString `gorm:"type:varchar(39)" db:"last_ipaddr"`
	IsManager      bool           `gorm:"type:tinyint(1);default:0;not null" db:"is_manager"`
	CDate          sql.NullString `gorm:"type:varchar(14);index" db:"cdate"`
	MDate          sql.NullString `gorm:"type:varchar(14);index" db:"mdate"`
}

func (Users) TableName() string {
	return "users"
}

var UsersColumns = fmt.Sprintf(" %s ", strings.Join(db.ColumnsForStruct(Users{}), ", "))
