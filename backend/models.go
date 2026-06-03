package main

import (
	"time"

	"github.com/google/uuid"
)

// Base model for UUIDs
type Base struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt time.Time
}

type Tenant struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string    `gorm:"size:255;not null"`
	Domain    string    `gorm:"size:255;uniqueIndex;not null"` // e.g. pbx.company.com
	CreatedAt time.Time
	UpdatedAt time.Time
}

// User represents a dashboard admin for a specific tenant
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TenantID     uuid.UUID `gorm:"type:uuid;index;not null"`
	Email        string    `gorm:"size:255;uniqueIndex;not null"`
	PasswordHash string    `gorm:"not null"`
	Role         string    `gorm:"size:50;default:'admin'"` // 'superadmin', 'admin', 'viewer'
	CreatedAt    time.Time
}

type Extension struct {
	Base
	TenantID         uuid.UUID `gorm:"type:uuid;not null"`
	ExtensionNumber  string    `gorm:"size:50;not null;uniqueIndex:idx_tenant_ext"`
	PasswordHash     string    `gorm:"size:255;not null"`
	IsActive         bool      `gorm:"default:true"`
	ProvisioningCode string    `gorm:"size:8;uniqueIndex"`
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
	TenantID       uuid.UUID `gorm:"type:uuid"` // Optional for system calls, but links to tenant
	CallerIDName   string `gorm:"size:100"`
	CallerIDNumber string `gorm:"size:50"`
	Destination    string `gorm:"size:100"`
	Context        string `gorm:"size:50"`
	Duration       int    `gorm:"default:0"` // Total duration
	Billsec        int    `gorm:"default:0"` // Answered duration
	HangupCause    string `gorm:"size:100"`
	UUID           string `gorm:"size:100;uniqueIndex"` // Call UUID
	Transcription  string `gorm:"type:text"` // Stores AI text dictation/transcription
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

// Voicemail stores recorded messages
type Voicemail struct {
	Base
	TenantID  uuid.UUID `gorm:"type:uuid;not null"`
	Extension string    `gorm:"size:50;not null"`
	CallerID  string    `gorm:"size:100"`
	Duration  int       `gorm:"default:0"`
	FilePath  string    `gorm:"size:255;not null"`
	IsRead    bool      `gorm:"default:false"`
}

// VideoRoom defines a WebRTC video conferencing room
type VideoRoom struct {
	Base
	TenantID uuid.UUID `gorm:"type:uuid;not null"`
	Name     string    `gorm:"size:100;not null"`
	RoomCode string    `gorm:"size:50;uniqueIndex;not null"`
	MaxUsers int       `gorm:"default:10"`
	IsActive bool      `gorm:"default:true"`
}

// IPCamera defines an endpoint for surveillance video streams
type IPCamera struct {
	Base
	TenantID uuid.UUID `gorm:"type:uuid;not null"`
	Name     string    `gorm:"size:100;not null"`
	RTSPUrl  string    `gorm:"size:255;not null"`
	Location string    `gorm:"size:100"`
}

// TeamChat defines a messaging channel
type TeamChat struct {
	Base
	TenantID uuid.UUID `gorm:"type:uuid;not null"`
	Name     string    `gorm:"size:100;not null"`
	Type     string    `gorm:"size:50;not null"` // "public", "private"
}

// SMSMessage stores inbound/outbound text messages
type SMSMessage struct {
	Base
	TenantID  uuid.UUID `gorm:"type:uuid;not null"`
	Direction string    `gorm:"size:20;not null"` // "inbound", "outbound"
	Sender    string    `gorm:"size:50;not null"`
	Receiver  string    `gorm:"size:50;not null"`
	Body      string    `gorm:"type:text;not null"`
}

// APICredential stores external API access tokens
type APICredential struct {
	Base
	TenantID uuid.UUID `gorm:"type:uuid;not null"`
	Name     string    `gorm:"size:100;not null"`
	TokenKey string    `gorm:"size:100;uniqueIndex;not null"`
}

// Webhook defines an HTTP callback for system events
type Webhook struct {
	Base
	TenantID uuid.UUID `gorm:"type:uuid;not null"`
	Event    string    `gorm:"size:100;not null"`
	TargetUrl string   `gorm:"size:255;not null"`
}

// ThrottleRule limits API or SIP requests
type ThrottleRule struct {
	Base
	TenantID uuid.UUID `gorm:"type:uuid;not null"`
	Name     string    `gorm:"size:100;not null"`
	Rate     int       `gorm:"default:60"`
}

// FirewallRule manages iptables / network ACLs
type FirewallRule struct {
	Base
	TenantID   uuid.UUID `gorm:"type:uuid;not null"`
	IPAddress  string    `gorm:"size:50;not null"`
	Action     string    `gorm:"size:20;not null"` // "allow", "block"
	Notes      string    `gorm:"size:255"`
}

// FraudRule manages destination prefix blocking
type FraudRule struct {
	Base
	TenantID    uuid.UUID `gorm:"type:uuid;not null"`
	Prefix      string    `gorm:"size:50;not null"`
	MaxCost     float64   `gorm:"default:0"`
	BlockAction string    `gorm:"size:50;not null"`
}
