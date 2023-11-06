package vbase

import (
	"database/sql"
	"fmt"
	"restApi/util/db"
	"strings"
)

type DeletedRoomOptions struct {
	RoomCode     string         `gorm:"type:varchar(12);not null;primaryKey" db:"room_code"`
	Name         string         `gorm:"type:varchar(128);not null;primaryKey" db:"name"`
	Item         string         `gorm:"type:varchar(64);not null;default:'';primaryKey" db:"item"`
	OptionType   int            `gorm:"type:int;not null" db:"option_type"` // 0:SERVER_ONLY, 1:CLIENT_ONLY, 2:BOTH
	ValueType    int            `gorm:"type:int;not null" db:"value_type"`  // 0:BOOLEAN, 1:NUMBER, 2:STRING, ...
	DefaultValue string         `gorm:"type:mediumtext;not null;default:''" db:"default_value"`
	OptionValue  string         `gorm:"type:mediumtext;not null;default:''" db:"option_value"`
	Notes        sql.NullString `gorm:"type:nvarchar(1024)" db:"notes"`
	CDate        sql.NullString `gorm:"type:varchar(14);index" db:"cdate"`
	MDate        sql.NullString `gorm:"type:varchar(14);index" db:"mdate"`
}

func (DeletedRoomOptions) TableName() string {
	return "deleted_room_options"
}

var DeletedRoomOptionsColumns = fmt.Sprintf(" %s ", strings.Join(db.ColumnsForStruct(DeletedRoomOptions{}), ", "))
