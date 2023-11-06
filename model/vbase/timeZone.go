package vbase

import (
	"database/sql"
	"fmt"
	"restApi/util/db"
	"strings"
)

type TimeZone struct {
	CountryCode  string         `gorm:"type:char(2);not null" db:"country_code"`
	TimeZone     string         `gorm:"type:varchar(32);not null;primaryKey" db:"time_zone"`
	UTCOffset    int            `gorm:"not null" db:"utc_offset"`
	UTCDSTOffset sql.NullInt16  `gorm:"" db:"utc_dst_offset"` // Nullable
	UseDST       int            `gorm:"default:0;not null" db:"use_dst"`
	Notes        sql.NullString `gorm:"type:varchar(256)" db:"notes"` // Nullable
}

func (TimeZone) TableName() string {
	return "time_zones"
}

var TimeZoneColumns = fmt.Sprintf(" %s ", strings.Join(db.ColumnsForStruct(TimeZone{}), ", "))
