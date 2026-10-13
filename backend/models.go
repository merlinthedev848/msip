package main

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base model for UUIDs
type Base struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"ID"`
	CreatedAt time.Time `json:"CreatedAt"`
}

// BeforeCreate guarantees a server-generated primary key
func (b *Base) BeforeCreate(_ *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

func (b *Base) GetID() uuid.UUID            { return b.ID }
func (b *Base) SetID(id uuid.UUID)          { b.ID = id }
func (b *Base) GetCreatedAt() time.Time     { return b.CreatedAt }
func (b *Base) SetCreatedAt(t time.Time)    { b.CreatedAt = t }

// TenantRef is embedded by every tenant-owned resource.
type TenantRef struct {
	TenantID uuid.UUID `gorm:"type:uuid;not null;index" json:"TenantID"`
}

func (t *TenantRef) GetTenantID() uuid.UUID   { return t.TenantID }
func (t *TenantRef) SetTenantID(id uuid.UUID) { t.TenantID = id }

type Tenant struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"ID"`
	Name      string    `gorm:"size:255;not null" json:"Name"`
	Domain    string    `gorm:"size:255;uniqueIndex;not null" json:"Domain"` // e.g. pbx.company.com
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
}

func (t *Tenant) BeforeCreate(_ *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// User represents a dashboard admin for a specific tenant
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"ID"`
	TenantID     uuid.UUID `gorm:"type:uuid;index;not null" json:"TenantID"`
	Email        string    `gorm:"size:255;uniqueIndex;not null" json:"Email"`
	PasswordHash string    `gorm:"not null" json:"PasswordHash,omitempty"` // excluded or sanitized in responses
	Role         string    `gorm:"size:50;default:'admin'" json:"Role"`     // 'superadmin', 'admin', 'viewer'
	CreatedAt    time.Time `json:"CreatedAt"`
}

func (u *User) BeforeCreate(_ *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

type Extension struct {
	Base
	TenantRef
	ExtensionNumber  string `gorm:"size:50;not null;uniqueIndex:idx_tenant_ext" json:"ExtensionNumber"`
	PasswordHash     string `gorm:"size:255;not null" json:"PasswordHash,omitempty"` // SIP secret/password
	IsActive         bool   `gorm:"default:true" json:"IsActive"`
	ProvisioningCode string `gorm:"size:8;uniqueIndex" json:"ProvisioningCode"`
}

type InboundRoute struct {
	Base
	TenantRef
	DIDNumber         string `gorm:"size:50;not null;uniqueIndex:idx_tenant_did" json:"DIDNumber"`
	DestinationType   string `gorm:"size:50;not null" json:"DestinationType"`
	DestinationTarget string `gorm:"size:100;not null" json:"DestinationTarget"`
}

type OutboundRoute struct {
	Base
	TenantRef
	DialPattern string    `gorm:"size:100;not null" json:"DialPattern"`
	TrunkID     uuid.UUID `gorm:"type:uuid;not null;index" json:"TrunkID"`
	Priority    int       `gorm:"default:1" json:"Priority"`
}

type Trunk struct {
	Base
	TenantRef
	ProviderName string `gorm:"size:100;not null" json:"ProviderName"`
	SIPServer    string `gorm:"size:255;not null" json:"SIPServer"`
	AuthMethod   string `gorm:"size:50;not null" json:"AuthMethod"` // IP_ACL or USERPASS
	Username     string `gorm:"size:100" json:"Username,omitempty"`
	Password     string `gorm:"size:100" json:"Password,omitempty"`
	Register     bool   `gorm:"default:false" json:"Register"`
}

// CDR stores Call Detail Records pushed by FreeSWITCH
type CDR struct {
	Base
	TenantID       uuid.UUID `gorm:"type:uuid;index" json:"TenantID"` // links to tenant
	CallerIDName   string    `gorm:"size:100" json:"CallerIDName"`
	CallerIDNumber string    `gorm:"size:50" json:"CallerIDNumber"`
	Destination    string    `gorm:"size:100" json:"Destination"`
	Context        string    `gorm:"size:50" json:"Context"`
	Duration       int       `gorm:"default:0" json:"Duration"`
	Billsec        int       `gorm:"default:0" json:"Billsec"`
	HangupCause    string    `gorm:"size:100" json:"HangupCause"`
	UUID           string    `gorm:"size:100;uniqueIndex" json:"UUID"`
	Transcription  string    `gorm:"type:text" json:"Transcription"`
}

// IVRMenu defines an Auto Attendant menu
type IVRMenu struct {
	Base
	TenantRef
	Name        string      `gorm:"size:100;not null" json:"Name"`
	Extension   string      `gorm:"size:50;not null;uniqueIndex:idx_tenant_ivr_ext" json:"Extension"`
	Greeting    string      `gorm:"size:255;not null" json:"Greeting"` // Path to wav or text for TTS
	Timeout     int         `gorm:"default:10000" json:"Timeout"`     // Milliseconds
	MaxFailures int         `gorm:"default:3" json:"MaxFailures"`
	Choices     []IVRChoice `gorm:"foreignKey:MenuID" json:"Choices"`
}

// IVRChoice defines digit mapping for a menu
type IVRChoice struct {
	Base
	MenuID      uuid.UUID `gorm:"type:uuid;not null;index" json:"MenuID"`
	Digits      string    `gorm:"size:10;not null" json:"Digits"`
	Action      string    `gorm:"size:50;not null" json:"Action"`
	Destination string    `gorm:"size:255;not null" json:"Destination"`
}

// Voicemail stores recorded messages
type Voicemail struct {
	Base
	TenantRef
	Extension string `gorm:"size:50;not null" json:"Extension"`
	CallerID  string `gorm:"size:100" json:"CallerID"`
	Duration  int    `gorm:"default:0" json:"Duration"`
	FilePath  string `gorm:"size:255;not null" json:"FilePath"`
	IsRead    bool   `gorm:"default:false" json:"IsRead"`
}

// VideoRoom defines a WebRTC video conferencing room
type VideoRoom struct {
	Base
	TenantRef
	Name     string `gorm:"size:100;not null" json:"Name"`
	RoomCode string `gorm:"size:50;uniqueIndex;not null" json:"RoomCode"`
	MaxUsers int    `gorm:"default:10" json:"MaxUsers"`
	IsActive bool   `gorm:"default:true" json:"IsActive"`
}

// IPCamera defines an endpoint for surveillance video streams
type IPCamera struct {
	Base
	TenantRef
	Name     string `gorm:"size:100;not null" json:"Name"`
	RTSPUrl  string `gorm:"size:255;not null" json:"RTSPUrl"`
	Location string `gorm:"size:100" json:"Location"`
}

// TeamChat defines a messaging channel
type TeamChat struct {
	Base
	TenantRef
	Name string `gorm:"size:100;not null" json:"Name"`
	Type string `gorm:"size:50;not null" json:"Type"` // "public", "private"
}

// SMSMessage stores inbound/outbound text messages
type SMSMessage struct {
	Base
	TenantRef
	Direction string `gorm:"size:20;not null" json:"Direction"` // "inbound", "outbound"
	Sender    string `gorm:"size:50;not null" json:"Sender"`
	Receiver  string `gorm:"size:50;not null" json:"Receiver"`
	Body      string `gorm:"type:text;not null" json:"Body"`
}

// APICredential stores external API access tokens
type APICredential struct {
	Base
	TenantRef
	Name       string `gorm:"size:100;not null" json:"Name"`
	TokenKey   string `gorm:"size:100;uniqueIndex;not null" json:"TokenKey"`
	SecretHash string `gorm:"size:64" json:"-"`
	Secret     string `gorm:"-" json:"Secret,omitempty"`
}

// Webhook defines an HTTP callback for system events
type Webhook struct {
	Base
	TenantRef
	Event     string `gorm:"size:100;not null" json:"Event"`
	TargetUrl string `gorm:"size:255;not null" json:"TargetUrl"`
}

// ThrottleRule limits API or SIP requests
type ThrottleRule struct {
	Base
	TenantRef
	Name string `gorm:"size:100;not null" json:"Name"`
	Rate int    `gorm:"default:60" json:"Rate"`
}

// FirewallRule manages iptables / network ACLs
type FirewallRule struct {
	Base
	TenantRef
	IPAddress string `gorm:"size:50;not null" json:"IPAddress"`
	Action    string `gorm:"size:20;not null" json:"Action"` // "allow"/"accept", "block"/"drop"
	Notes     string `gorm:"size:255" json:"Notes"`
}

// FraudRule manages destination prefix blocking
type FraudRule struct {
	Base
	TenantRef
	Prefix      string  `gorm:"size:50;not null" json:"Prefix"`
	MaxCost     float64 `gorm:"default:0" json:"MaxCost"`
	BlockAction string  `gorm:"size:50;not null" json:"BlockAction"`
}

// TimeCondition defines business hours routing
type TimeCondition struct {
	Base
	TenantRef
	Name              string `gorm:"size:100;not null" json:"Name"`
	DaysOfWeek        string `gorm:"size:100;default:'Mon,Tue,Wed,Thu,Fri'" json:"DaysOfWeek"`
	StartTime         string `gorm:"size:10;default:'09:00'" json:"StartTime"`
	EndTime           string `gorm:"size:10;default:'17:00'" json:"EndTime"`
	InOfficeDestType  string `gorm:"size:50;default:'EXTENSION'" json:"InOfficeDestType"`
	InOfficeTarget    string `gorm:"size:100;not null" json:"InOfficeTarget"`
	OutOfOfficeDestType string `gorm:"size:50;default:'VOICEMAIL'" json:"OutOfOfficeDestType"`
	OutOfOfficeTarget   string `gorm:"size:100;not null" json:"OutOfOfficeTarget"`
}

// Holiday defines scheduled holidays
type Holiday struct {
	Base
	TenantRef
	Name       string `gorm:"size:100;not null" json:"Name"`
	StartDate  string `gorm:"size:50;not null" json:"StartDate"`
	EndDate    string `gorm:"size:50;not null" json:"EndDate"`
	DestType   string `gorm:"size:50;default:'VOICEMAIL'" json:"DestType"`
	DestTarget string `gorm:"size:100;not null" json:"DestTarget"`
}

// OutboundRestriction blocks dialing specific prefixes
type OutboundRestriction struct {
	Base
	TenantRef
	Name        string `gorm:"size:100;not null" json:"Name"`
	Pattern     string `gorm:"size:100;not null" json:"Pattern"` // e.g. 0900, 00, +
	Extensions  string `gorm:"size:255;default:'ALL'" json:"Extensions"` // "ALL" or comma-separated exts
	Action      string `gorm:"size:50;default:'BLOCK'" json:"Action"` // BLOCK or PIN
}

// AutoCLIP routes returning callers directly to the extension that last phoned them
type AutoCLIP struct {
	Base
	TenantRef
	CallerNumber string    `gorm:"size:50;not null;index" json:"CallerNumber"`
	Extension    string    `gorm:"size:50;not null" json:"Extension"`
	ExpiresAt    time.Time `json:"ExpiresAt"`
}

// PagingGroup defines intercom/paging broadcast groups
type PagingGroup struct {
	Base
	TenantRef
	Name      string `gorm:"size:100;not null" json:"Name"`
	Extension string `gorm:"size:50;not null;uniqueIndex:idx_tenant_paging_ext" json:"Extension"`
	Members   string `gorm:"type:text;not null" json:"Members"` // comma-separated extensions
	TwoWay    bool   `gorm:"default:false" json:"TwoWay"` // false=1-way paging, true=2-way intercom
}

// PickupGroup allows users to grab ringing calls with *8
type PickupGroup struct {
	Base
	TenantRef
	Name    string `gorm:"size:100;not null" json:"Name"`
	Code    string `gorm:"size:50;not null" json:"Code"` // e.g. *8
	Members string `gorm:"type:text;not null" json:"Members"` // comma-separated extensions
}

// AudioPrompt manages MoH playlists, system prompts, and custom audio files
type AudioPrompt struct {
	Base
	TenantRef
	Name        string `gorm:"size:100;not null" json:"Name"`
	Category    string `gorm:"size:50;not null" json:"Category"` // "MOH", "SYSTEM", "CUSTOM"
	FileName    string `gorm:"size:255;not null" json:"FileName"`
	Description string `gorm:"size:255" json:"Description"`
}

// FeatureCode lists and configures dialed feature codes
type FeatureCode struct {
	Base
	TenantRef
	Name        string `gorm:"size:100;not null" json:"Name"`
	Code        string `gorm:"size:20;not null" json:"Code"`
	DefaultCode string `gorm:"size:20;not null" json:"DefaultCode"`
	Description string `gorm:"size:255" json:"Description"`
	IsActive    bool   `gorm:"default:true" json:"IsActive"`
}

// RecordingRule configures automated call recording
type RecordingRule struct {
	Base
	TenantRef
	Scope       string `gorm:"size:50;not null" json:"Scope"` // "GLOBAL", "INBOUND", "OUTBOUND", "EXTENSION"
	Target      string `gorm:"size:100;default:'ALL'" json:"Target"`
	RecordMode  string `gorm:"size:50;default:'ALWAYS'" json:"RecordMode"` // "ALWAYS", "ON_DEMAND"
	StorageDays int    `gorm:"default:90" json:"StorageDays"`
}
