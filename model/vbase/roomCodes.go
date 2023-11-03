package vbase

import (
	"fmt"
	"restApi/util/db"
	"strings"
)

type RoomCodes struct {
	Idx      string `gorm:"type:char(36);primaryKey" db:"idx"`
	RoomCode string `gorm:"type:varchar(12);not null;unique" db:"room_code"`
	RoomID   string `gorm:"type:varchar(36)" db:"room_id"`
	CDate    string `gorm:"type:varchar(14);index" db:"cdate"`
	MDate    string `gorm:"type:varchar(14);index" db:"mdate"`
}

func (RoomCodes) TableName() string {
	return "room_codes"
}

var RoomCodesColumns = fmt.Sprintf(" %s ", strings.Join(db.ColumnsForStruct(RoomCodes{}), ", "))
