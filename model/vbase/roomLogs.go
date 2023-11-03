package vbase

import (
	"fmt"
	"restApi/util/db"
	"strings"
)

type RoomLogs struct {
	SiteIdx       string `gorm:"type:char(36);not null;index" db:"site_idx"`
	UserIdx       string `gorm:"type:char(36);not null;index" db:"user_idx"`
	UserID        string `gorm:"type:varchar(128);not null;index" db:"user_id"`
	RoomID        string `gorm:"type:char(36);not null;index" db:"room_id"`
	RoomCode      string `gorm:"type:varchar(12);not null;index" db:"room_code"`
	RoomPolicy    string `gorm:"type:varchar(16)" db:"room_policy"`
	ServerIdx     string `gorm:"type:char(36);index" db:"server_idx"`
	Passwd        string `gorm:"type:varchar(32)" db:"passwd"`
	Title         string `gorm:"type:nvarchar(128);not null" db:"title"`
	TimeZone      string `gorm:"type:varchar(32);not null" db:"time_zone"`
	IsPublic      int8   `gorm:"type:tinyint;not null" db:"is_public"`
	MaxUsers      int    `gorm:"type:int;not null" db:"max_users"`
	AdmissionDate string `gorm:"type:varchar(14);index" db:"admission_date"`
	PlannedDate   string `gorm:"type:varchar(14);not null;index" db:"planned_date"`
	RoomDuration  int    `gorm:"type:int;not null" db:"room_duration"`
	CDate         string `gorm:"type:varchar(14);index" db:"cdate"`
	MDate         string `gorm:"type:varchar(14);index" db:"mdate"`
	// Added for log
	Idx          string `gorm:"type:char(36);primaryKey" db:"idx"`
	InstanceIdx  string `gorm:"type:char(36);not null;index" db:"instance_idx"`
	StartedDate  string `gorm:"type:varchar(14);not null;index" db:"started_date"`
	FinishedDate string `gorm:"type:varchar(14);index" db:"finished_date"`
	LDate        string `gorm:"type:varchar(14);index" db:"ldate"`
}

func (RoomLogs) TableName() string {
	return "room_logs"
}

var RoomLogsColumns = fmt.Sprintf(" %s ", strings.Join(db.ColumnsForStruct(RoomLogs{}), ", "))
