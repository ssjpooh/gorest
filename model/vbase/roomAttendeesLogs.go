package vbase

import (
	"fmt"
	"restApi/util/db"
	"strings"
)

type RoomAttendeeLogs struct {
	SiteIdx       string `gorm:"type:char(36);not null;index" db:"site_idx"`
	UserIdx       string `gorm:"type:char(36);index" db:"user_idx"`
	UserID        string `gorm:"type:varchar(128);index" db:"user_id"`
	RoomCode      string `gorm:"type:varchar(12);not null;index" db:"room_code"`
	AttdID        string `gorm:"type:varchar(128);not null;index" db:"attd_id"`
	IsManager     int8   `gorm:"type:tinyint;default:0;not null" db:"is_manager"`
	IsSubManager  int8   `gorm:"type:tinyint;default:0;not null" db:"is_sub_manager"`
	Email         string `gorm:"type:varchar(128);default:''" db:"email"`
	Name          string `gorm:"type:nvarchar(64);not null" db:"name"`
	NickName      string `gorm:"type:nvarchar(64)" db:"nick_name"`
	InitialRights string `gorm:"type:varchar(128);default:'';not null" db:"initial_rights"`
	AttendedDate  string `gorm:"type:varchar(14);not null;index" db:"attended_date"`
	ExitedDate    string `gorm:"type:varchar(14);index" db:"exited_date"`
	IPAddr        string `gorm:"type:varchar(39)" db:"ipaddr"`
	CDate         string `gorm:"type:varchar(14);index" db:"cdate"`
	MDate         string `gorm:"type:varchar(14);index" db:"mdate"`
	// Log specific fields
	Idx         string `gorm:"type:char(36);not null;primaryKey" db:"idx"`
	ServerIdx   string `gorm:"type:char(36);not null;index" db:"server_idx"`
	InstanceIdx string `gorm:"type:char(36);not null;index" db:"instance_idx"`
	LDate       string `gorm:"type:varchar(14);index" db:"ldate"`
}

func (RoomAttendeeLogs) TableName() string {
	return "room_attendee_logs"
}

var RoomAttendeeLogsColumns = fmt.Sprintf(" %s ", strings.Join(db.ColumnsForStruct(RoomAttendeeLogs{}), ", "))
