package vbase

import (
	"database/sql"
	"fmt"
	"restApi/util/db"
	"strings"
)

type Rooms struct {
	SiteIdx       string         `gorm:"type:char(36);not null;index" db:"site_idx"`
	UserIdx       string         `gorm:"type:char(36);not null;index" db:"user_idx"`
	UserID        string         `gorm:"type:varchar(128);not null;index" db:"user_id"`
	RoomID        string         `gorm:"type:char(36);primaryKey" db:"room_id"`
	RoomCode      string         `gorm:"type:varchar(12);uniqueIndex" db:"room_code"` // Though it's nullable, the unique constraint still applies
	RoomPolicy    sql.NullString `gorm:"type:varchar(16)" db:"room_policy"`
	ServerIdx     sql.NullString `gorm:"type:char(36);index" db:"server_idx"`
	Passwd        sql.NullString `gorm:"type:varchar(32)" db:"passwd"`
	Title         string         `gorm:"type:nvarchar(128);not null" db:"title"`
	TimeZone      sql.NullString `gorm:"type:varchar(32);default:'Asia/Seoul';not null" db:"time_zone"`
	IsPublic      int8           `gorm:"type:tinyint;default:1;not null" db:"is_public"`
	MaxUsers      int            `gorm:"type:int;default:-1;not null" db:"max_users"`
	AdmissionDate sql.NullString `gorm:"type:varchar(14);index" db:"admission_date"`
	PlannedDate   sql.NullString `gorm:"type:varchar(14);not null;index" db:"planned_date"`
	RoomDuration  int            `gorm:"type:int;default:-1;not null" db:"room_duration"`
	CDate         sql.NullString `gorm:"type:varchar(14);index" db:"cdate"`
	MDate         sql.NullString `gorm:"type:varchar(14);index" db:"mdate"`
}

func (Rooms) TableName() string {
	return "rooms"
}

var RoomsColumns = fmt.Sprintf(" %s ", strings.Join(db.ColumnsForStruct(Rooms{}), ", "))
