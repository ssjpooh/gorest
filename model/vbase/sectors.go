package vbase

import (
	"fmt"
	"restApi/util/db"
	"strings"
)

type Sectors struct {
	Sector string `gorm:"type:varchar(16);primaryKey" db:"sector"`
	Notes  string `gorm:"type:nvarchar(1024)" db:"notes"`
	CDate  string `gorm:"type:varchar(14);index" db:"cdate"`
}

func (Sectors) TableName() string {
	return "sectors"
}

var SectorColumns = fmt.Sprintf(" %s ", strings.Join(db.ColumnsForStruct(Sectors{}), ", "))
