package vbase

import (
	"fmt"
	"restApi/util/db"
	"strings"
)

type Servers struct {
	Sector         string `gorm:"type:varchar(16);not null" db:"sector"`
	ServerIdx      string `gorm:"type:char(36);primaryKey" db:"server_idx"`
	ServerTypes    string `gorm:"type:varchar(128)" db:"server_types"`
	IsRegistered   int    `gorm:"type:tinyint;default:0;not null" db:"is_registered"`
	IsActivated    int    `gorm:"type:tinyint;default:0;not null" db:"is_activated"`
	IsAllowed      int    `gorm:"type:tinyint;default:1;not null" db:"is_allowed"`
	PrivateIpaddrs string `gorm:"type:varchar(200)" db:"private_ipaddrs"`
	PublicIpaddr   string `gorm:"type:varchar(39)" db:"public_ipaddr"`
	PublicDomains  string `gorm:"type:varchar(256)" db:"public_domains"`
	Version        string `gorm:"type:varchar(16)" db:"version"`
	StartedDate    string `gorm:"type:varchar(14)" db:"started_date"`
	StoppedDate    string `gorm:"type:varchar(14)" db:"stopped_date"`
	GorTotal       int    `gorm:"type:int;default:0;not null" db:"gor_total"`
	Gor            int    `gorm:"type:int;default:0;not null" db:"gor"`
	CpuTotal       int    `gorm:"type:int;default:0;not null" db:"cpu_total"`
	Cpu            int    `gorm:"type:int;default:0;not null" db:"cpu"`
	MemTotal       int    `gorm:"type:int;default:0;not null" db:"mem_total"`
	Mem            int    `gorm:"type:int;default:0;not null" db:"mem"`
	HddTotal       int    `gorm:"type:int;default:0;not null" db:"hdd_total"`
	Hdd            int    `gorm:"type:int;default:0;not null" db:"hdd"`
	RtpTotal       int    `gorm:"type:int;default:0;not null" db:"rtp_total"`
	Rtp            int    `gorm:"type:int;default:0;not null" db:"rtp"`
	Dumps          int    `gorm:"type:int;default:0;not null" db:"dumps"`
	Rooms          int    `gorm:"type:int;default:0;not null" db:"rooms"`
	Attendees      int    `gorm:"type:int;default:0;not null" db:"attendees"`
	Notes          string `gorm:"type:nvarchar(1024)" db:"notes"`
	CDate          string `gorm:"type:varchar(14)" db:"cdate"`
	MDate          string `gorm:"type:varchar(14)" db:"mdate"`
}

func (Servers) TableName() string {
	return "servers"
}

var ServerColumns = fmt.Sprintf(" %s ", strings.Join(db.ColumnsForStruct(Servers{}), ", "))
