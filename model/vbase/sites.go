package vbase

import (
	"fmt"
	"restApi/util/db"
	"strings"
)

type Sites struct {
	Sector         string `gorm:"type:varchar(16);not null" db:"sector"`
	SiteIdx        string `gorm:"type:char(36);primaryKey" db:"site_idx"`
	SiteID         string `gorm:"type:varchar(32);unique;not null" db:"site_id"`
	Passwd         string `gorm:"type:varchar(32)" db:"passwd"`
	Email          string `gorm:"type:varchar(128);default:''" db:"email"`
	Name           string `gorm:"type:nvarchar(64);not null" db:"name"`
	CompanyName    string `gorm:"type:varchar(64)" db:"company_name"`
	LastLogoutDate string `gorm:"type:varchar(14)" db:"last_logout_date"`
	CDate          string `gorm:"type:varchar(14)" db:"cdate"`
	MDate          string `gorm:"type:varchar(14)" db:"mdate"`
}

func (Sites) TableName() string {
	return "sites"
}

var SitesColumns = fmt.Sprintf(" %s ", strings.Join(db.ColumnsForStruct(Sites{}), ", "))
