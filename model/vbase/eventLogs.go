package vbase

import (
	"fmt"
	"restApi/util/db"
	"strings"
)

type EventLogs struct {
	SiteIdx       *string `gorm:"type:char(36);index" db:"site_idx"`
	UserIdx       *string `gorm:"type:char(36);index" db:"user_idx"`
	UserID        *string `gorm:"type:varchar(128);index" db:"user_id"`
	RoomCode      *string `gorm:"type:varchar(12);index" db:"room_code"`
	AttdID        *string `gorm:"type:varchar(128);index" db:"attd_id"`
	IsManager     *int    `gorm:"type:tinyint" db:"is_manager"`
	IsSubManager  *int    `gorm:"type:tinyint" db:"is_sub_manager"`
	Email         *string `gorm:"type:varchar(128);index" db:"email"`
	Name          *string `gorm:"type:nvarchar(64);index" db:"name"`
	NickName      *string `gorm:"type:nvarchar(64);index" db:"nick_name"`
	InitialRights *string `gorm:"type:varchar(128)" db:"initial_rights"`
	AttendedDate  *string `gorm:"type:varchar(14);index" db:"attended_date"`
	ExitedDate    *string `gorm:"type:varchar(14);index" db:"exited_date"`
	IPAddr        *string `gorm:"type:varchar(39);index" db:"ipaddr"`
	CDate         *string `gorm:"type:varchar(14);index" db:"cdate"`
	MDate         *string `gorm:"type:varchar(14);index" db:"mdate"`
	Idx           string  `gorm:"type:char(36);not null;primaryKey" db:"idx"`
	EventCode     string  `gorm:"type:varchar(32);not null;index" db:"event_code"`
	EventContent  *string `gorm:"type:mediumtext" db:"event_content"`
	ServerIdx     *string `gorm:"type:char(36);index" db:"server_idx"`
	InstanceIdx   *string `gorm:"type:char(36);index" db:"instance_idx"`
	LDate         string  `gorm:"type:varchar(14);not null;index" db:"ldate"`
}

func (EventLogs) TableName() string {
	return "event_logs"
}

var EventLogsColumns = fmt.Sprintf(" %s ", strings.Join(db.ColumnsForStruct(EventLogs{}), ", "))
