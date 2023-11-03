package vbase

import (
	"fmt"
	"restApi/util/db"
	"strings"
)

type DeletedRoomChats struct {
	SiteIdx     string  `gorm:"type:char(36);not null;index" db:"site_idx"`
	RoomCode    string  `gorm:"type:varchar(12);not null;index" db:"room_code"`
	Idx         string  `gorm:"type:char(36);not null;primaryKey" db:"idx"`
	FilePath    string  `gorm:"type:nvarchar(300);not null" db:"file_path"`
	FileSize    int     `gorm:"type:int;default:0;not null" db:"file_size"`
	ServerIdx   *string `gorm:"type:char(36);index" db:"server_idx"`
	InstanceIdx *string `gorm:"type:char(36);index" db:"instance_idx"`
	CDate       *string `gorm:"type:varchar(14);index" db:"cdate"`
}

func (DeletedRoomChats) TableName() string {
	return "deleted_room_chats"
}

var DeletedRoomChatsColumns = fmt.Sprintf(" %s ", strings.Join(db.ColumnsForStruct(DeletedRoomChats{}), ", "))
