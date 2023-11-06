package vbase

import (
	"database/sql"
	"fmt"
	"restApi/util/db"
	"strings"
)

type EventLogs struct {
	SiteIdx       sql.NullString `gorm:"type:char(36);index" db:"site_idx"`
	UserIdx       sql.NullString `gorm:"type:char(36);index" db:"user_idx"`
	UserID        sql.NullString `gorm:"type:varchar(128);index" db:"user_id"`
	RoomCode      sql.NullString `gorm:"type:varchar(12);index" db:"room_code"`
	AttdID        sql.NullString `gorm:"type:varchar(128);index" db:"attd_id"`
	IsManager     sql.NullInt16  `gorm:"type:tinyint" db:"is_manager"`
	IsSubManager  sql.NullInt16  `gorm:"type:tinyint" db:"is_sub_manager"`
	Email         sql.NullString `gorm:"type:varchar(128);index" db:"email"`
	Name          sql.NullString `gorm:"type:nvarchar(64);index" db:"name"`
	NickName      sql.NullString `gorm:"type:nvarchar(64);index" db:"nick_name"`
	InitialRights sql.NullString `gorm:"type:varchar(128)" db:"initial_rights"`
	AttendedDate  sql.NullString `gorm:"type:varchar(14);index" db:"attended_date"`
	ExitedDate    sql.NullString `gorm:"type:varchar(14);index" db:"exited_date"`
	IPAddr        sql.NullString `gorm:"type:varchar(39);index" db:"ipaddr"`
	CDate         sql.NullString `gorm:"type:varchar(14);index" db:"cdate"`
	MDate         sql.NullString `gorm:"type:varchar(14);index" db:"mdate"`
	Idx           string         `gorm:"type:char(36);not null;primaryKey" db:"idx"`
	EventCode     string         `gorm:"type:varchar(32);not null;index" db:"event_code"`
	EventContent  sql.NullString `gorm:"type:mediumtext" db:"event_content"`
	ServerIdx     sql.NullString `gorm:"type:char(36);index" db:"server_idx"`
	InstanceIdx   sql.NullString `gorm:"type:char(36);index" db:"instance_idx"`
	LDate         string         `gorm:"type:varchar(14);not null;index" db:"ldate"`
}

func (EventLogs) TableName() string {
	return "event_logs"
}

var EventLogsColumns = fmt.Sprintf(" %s ", strings.Join(db.ColumnsForStruct(EventLogs{}), ", "))
