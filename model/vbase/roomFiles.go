package vbase

import (
	"fmt"
	"restApi/util/db"
	"strings"
)

type RoomFiles struct {
	SiteIdx   string `gorm:"type:char(36);not null;index" db:"site_idx"`
	UserIdx   string `gorm:"type:char(36);index" db:"user_idx"`
	UserID    string `gorm:"type:varchar(128);index" db:"user_id"`
	AttdID    string `gorm:"type:varchar(128);not null;index" db:"attd_id"`
	RoomCode  string `gorm:"type:varchar(12);not null;index" db:"room_code"`
	RoomGroup int    `gorm:"type:int;default:1;not null" db:"room_group"`
	FileIdx   string `gorm:"type:char(36);not null;primaryKey" db:"file_idx"`
	FileKind  string `gorm:"type:varchar(16);not null" db:"file_kind"`
	FileName  string `gorm:"type:nvarchar(128);not null" db:"file_name"`
	FilePath  string `gorm:"type:nvarchar(300);not null" db:"file_path"`
	FileSize  int    `gorm:"type:int;default:0;not null" db:"file_size"`
	Title     string `gorm:"type:nvarchar(128);not null" db:"title"`
	Pages     *int   `gorm:"type:int" db:"pages"`
	CDate     string `gorm:"type:varchar(14);index" db:"cdate"`
	MDate     string `gorm:"type:varchar(14);index" db:"mdate"`
}

func (RoomFiles) TableName() string {
	return "room_files"
}

var RoomFilesColumns = fmt.Sprintf(" %s ", strings.Join(db.ColumnsForStruct(RoomFiles{}), ", "))
