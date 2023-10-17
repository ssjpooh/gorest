package tenant

import (
	"time"

	"gorm.io/gorm"
)

type USERS struct {
	BusiCode                string  `gorm:"type:varchar(16);not null"`
	IDX                     string  `gorm:"type:char(36);primaryKey"`
	SeqNo                   int     `gorm:"autoIncrement;unique"`
	UserID                  string  `gorm:"type:varchar(128);not null;unique"`
	UserPassword            string  `gorm:"type:varchar(256);not null"`
	Email                   string  `gorm:"type:varchar(256);default:''"`
	CID                     string  `gorm:"type:char(13);not null"`
	Gender                  int     `gorm:"not null"`
	Birthday                *string `gorm:"type:char(8)"`
	IsLunar                 *int8
	DefUserName             string  `gorm:"type:varchar(256);not null"`
	EngUserName             string  `gorm:"type:varchar(256);not null"`
	DefCompanyName          *string `gorm:"type:varchar(256)"`
	EngCompanyName          *string `gorm:"type:varchar(256)"`
	CompanyPostNo           *string `gorm:"type:varchar(6)"`
	DefCompanyAddr1         *string `gorm:"type:varchar(64)"`
	EngCompanyAddr1         *string `gorm:"type:varchar(64)"`
	DefCompanyAddr2         *string `gorm:"type:varchar(64)"`
	EngCompanyAddr2         *string `gorm:"type:varchar(64)"`
	HomePostNo              *string `gorm:"type:varchar(6)"`
	DefHomeAddr1            *string `gorm:"type:varchar(64)"`
	EngHomeAddr1            *string `gorm:"type:varchar(64)"`
	DefHomeAddr2            *string `gorm:"type:varchar(64)"`
	EngHomeAddr2            *string `gorm:"type:varchar(64)"`
	HomePage                *string `gorm:"type:varchar(32)"`
	LocalCode               *int8
	AuthPhoneFlags          int     `gorm:"default:0;not null"`
	AuthPhoneSelected       int     `gorm:"default:0;not null"`
	AuthPhoneMobile         *string `gorm:"type:varchar(20)"`
	AuthPhoneHome           *string `gorm:"type:varchar(20)"`
	LastLogoutTime          *time.Time
	IsJoined                *int
	JoinDate                *time.Time
	CancelDate              *time.Time
	ProductName             *string `gorm:"type:varchar(64)"`
	PhoneKN                 *string `gorm:"type:varchar(20)"`
	DscTime                 *int
	DscPoint                *int
	ScTime                  *int
	ScPoint                 *int
	Remark                  *string `gorm:"type:varchar(1000)"`
	StartCallNo             int     `gorm:"default:100;not null"`
	CurrentCallNo           int     `gorm:"default:100;not null"`
	CompanyCode             *string `gorm:"type:varchar(256)"`
	SortNo                  *string `gorm:"type:varchar(256)"`
	IsUsed                  int8    `gorm:"default:1;not null"`
	LoginSecurityProcName   *string `gorm:"type:varchar(256)"`
	RefOrgIdx               *string `gorm:"type:char(36)"`
	TenantID                *string `gorm:"type:varchar(128)"`
	IsSNS                   *int8
	SNSCode                 *string `gorm:"type:varchar(1)"`
	SNSPid                  *string `gorm:"type:varchar(100)"`
	SNSImageURL             *string `gorm:"type:varchar(100)"`
	IsResetPW               *int8
	LastLoginTime           *time.Time
	IsEmailCertify          *int8
	StateCode               string `gorm:"type:varchar(2);default:'RG';not null"`
	LastPasswordChangedTime *time.Time
	LastLoginFailureTime    *time.Time
	WrongPasswordCounter    int     `gorm:"default:0;not null"`
	RemoveConfContext       int8    `gorm:"default:0;not null"`
	Category                *string `gorm:"type:varchar(16)"`
	UseWaitingRoom          int8    `gorm:"default:0;not null"`
	UseHostID               *int8
	HostIDCount             *int
	HostRoomUserCount       *int
	MaximumVolume           *int
	OneTimeVolume           *int
	RetentionPeriod         *int
	RecordingTotalSize      *int
	ManualRegUser           int8    `gorm:"default:0;not null"`
	VendorCertCode          *string `gorm:"type:varchar(28)"`
	RequestIPAddr           *string `gorm:"type:varchar(20)"`
	RequestSignUpDate       *time.Time
	StateCategory           string `gorm:"type:varchar(2);default:'DE';not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
