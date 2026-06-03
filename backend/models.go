package main

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base model for UUIDs
type Base struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt time.Time
}

type Tenant struct {
	Base
	Name string `gorm:"size:255;not null"`
}

type Extension struct {
	Base
	TenantID       uuid.UUID `gorm:"type:uuid;not null"`
	ExtensionNumber string    `gorm:"size:50;not null;uniqueIndex:idx_tenant_ext"`
	PasswordHash   string    `gorm:"size:255;not null"`
	IsActive       bool      `gorm:"default:true"`
}

type InboundRoute struct {
	Base
	TenantID          uuid.UUID `gorm:"type:uuid;not null"`
	DIDNumber         string    `gorm:"size:50;not null;uniqueIndex:idx_tenant_did"`
	DestinationType   string    `gorm:"size:50;not null"`
	DestinationTarget string    `gorm:"size:100;not null"`
}

type OutboundRoute struct {
	Base
	TenantID    uuid.UUID `gorm:"type:uuid;not null"`
	DialPattern string    `gorm:"size:100;not null"`
	TrunkID     uuid.UUID `gorm:"type:uuid;not null"`
	Priority    int       `gorm:"default:1"`
}

type Trunk struct {
	Base
	TenantID     uuid.UUID `gorm:"type:uuid;not null"`
	ProviderName string    `gorm:"size:100;not null"`
	SIPServer    string    `gorm:"size:255;not null"`
	AuthMethod   string    `gorm:"size:50;not null"` // IP_ACL or USERPASS
}

// CDR stores Call Detail Records pushed by FreeSWITCH
type CDR struct {
	Base
	CallerIDName   string `gorm:"size:100"`
	CallerIDNumber string `gorm:"size:50"`
	Destination    string `gorm:"size:100"`
	Context        string `gorm:"size:50"`
	Duration       int    `gorm:"default:0"` // Total duration
	Billsec        int    `gorm:"default:0"` // Answered duration
	HangupCause    string `gorm:"size:100"`
	UUID           string `gorm:"size:100;uniqueIndex"` // Call UUID
}

// IVRMenu defines an Auto Attendant menu
type IVRMenu struct {
	Base
	TenantID    uuid.UUID   `gorm:"type:uuid;not null"`
	Name        string      `gorm:"size:100;not null"`
	Extension   string      `gorm:"size:50;not null;uniqueIndex:idx_tenant_ivr_ext"`
	Greeting    string      `gorm:"size:255;not null"` // Path to wav or text for TTS
	Timeout     int         `gorm:"default:10000"` // Milliseconds
	MaxFailures int         `gorm:"default:3"`
	Choices     []IVRChoice `gorm:"foreignKey:MenuID"`
}

// IVRChoice defines digit mapping for a menu
type IVRChoice struct {
	Base
	MenuID      uuid.UUID `gorm:"type:uuid;not null"`
	Digits      string    `gorm:"size:10;not null"` // e.g. "1" or "0"
	Action      string    `gorm:"size:50;not null"` // e.g. "menu-exec-app"
	Destination string    `gorm:"size:255;not null"` // e.g. "transfer 1001 XML default"
}
