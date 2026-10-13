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

	// User Profile
	FirstName   string `gorm:"size:100" json:"FirstName"`
	LastName    string `gorm:"size:100" json:"LastName"`
	DisplayName string `gorm:"size:150" json:"DisplayName"`
	Email       string `gorm:"size:255" json:"Email"`
	Department  string `gorm:"size:100" json:"Department"`
	HomeNumber  string `gorm:"size:50" json:"HomeNumber"`

	// Telephony & Numbers
	CallerIDName          string     `gorm:"size:100" json:"CallerIDName"`
	CallerIDNumber        string     `gorm:"size:50" json:"CallerIDNumber"`
	OutboundDID           string     `gorm:"size:50" json:"OutboundDID"`
	DirectNumber          string     `gorm:"size:50" json:"DirectNumber"`
	PresentationNumber    string     `gorm:"size:50" json:"PresentationNumber"`
	SiteID                *uuid.UUID `gorm:"type:uuid" json:"SiteID"`
	ButtonProfileID       *uuid.UUID `gorm:"type:uuid" json:"ButtonProfileID"`
	VoicemailNotification string     `gorm:"size:50;default:'Off'" json:"VoicemailNotification"`

	// Forwarding & Presence
	MobileNumber           string `gorm:"size:50" json:"MobileNumber"`
	CallWaiting            bool   `gorm:"default:false" json:"CallWaiting"`
	DoNotDisturb           bool   `gorm:"default:false" json:"DoNotDisturb"`
	Forwarding             bool   `gorm:"default:false" json:"Forwarding"`
	ForwardInternalTarget  string `gorm:"size:100" json:"ForwardInternalTarget"`
	ForwardExternalTarget  string `gorm:"size:100" json:"ForwardExternalTarget"`
	MobileTwinning         bool   `gorm:"default:false" json:"MobileTwinning"`
	ForwardAlwaysEnabled   bool   `gorm:"default:false" json:"ForwardAlwaysEnabled"`
	ForwardAlwaysTarget    string `gorm:"size:100" json:"ForwardAlwaysTarget"`
	ForwardBusyEnabled     bool   `gorm:"default:false" json:"ForwardBusyEnabled"`
	ForwardBusyTarget      string `gorm:"size:100" json:"ForwardBusyTarget"`
	ForwardNoAnswerEnabled bool   `gorm:"default:false" json:"ForwardNoAnswerEnabled"`
	ForwardNoAnswerTarget  string `gorm:"size:100" json:"ForwardNoAnswerTarget"`
	RingTimeoutSeconds     int    `gorm:"default:20" json:"RingTimeoutSeconds"`

	// Internal Call Settings
	InternalNoAnswerTimeout int    `gorm:"default:3" json:"InternalNoAnswerTimeout"`
	InternalNoAnswerAction  string `gorm:"size:50;default:'Not Set'" json:"InternalNoAnswerAction"`
	InternalNoAnswerTarget  string `gorm:"size:100" json:"InternalNoAnswerTarget"`
	InternalOnBusyAction    string `gorm:"size:50;default:'Unknown'" json:"InternalOnBusyAction"`
	InternalOnBusyTarget    string `gorm:"size:100" json:"InternalOnBusyTarget"`

	// External Call Settings
	ExternalNoAnswerTimeout int    `gorm:"default:0" json:"ExternalNoAnswerTimeout"`
	ExternalNoAnswerAction  string `gorm:"size:50;default:'Not Set'" json:"ExternalNoAnswerAction"`
	ExternalNoAnswerTarget  string `gorm:"size:100" json:"ExternalNoAnswerTarget"`
	ExternalOnBusyAction    string `gorm:"size:50;default:'Unknown'" json:"ExternalOnBusyAction"`
	ExternalOnBusyTarget    string `gorm:"size:100" json:"ExternalOnBusyTarget"`

	// Handset Programmable Buttons (JSON string)
	HandsetButtons string `gorm:"type:text;default:'[]'" json:"HandsetButtons"`

	// Voicemail Preferences
	VoicemailEnabled bool   `gorm:"default:true" json:"VoicemailEnabled"`
	VoicemailPIN     string `gorm:"size:20" json:"VoicemailPIN"`
	VoicemailToEmail bool   `gorm:"default:true" json:"VoicemailToEmail"`
}

type InboundRoute struct {
	Base
	TenantRef
	DIDNumber           string `gorm:"size:50;not null;uniqueIndex:idx_tenant_did" json:"DIDNumber"`
	CallerIDPattern     string `gorm:"size:100" json:"CallerIDPattern"`
	Name                string `gorm:"size:100" json:"Name"`
	DestinationType     string `gorm:"size:50;not null" json:"DestinationType"`
	DestinationTarget   string `gorm:"size:100;not null" json:"DestinationTarget"`
	FallbackDestType    string `gorm:"size:50" json:"FallbackDestType"`
	FallbackDestTarget  string `gorm:"size:100" json:"FallbackDestTarget"`
	TimeConditionID     string `gorm:"size:100" json:"TimeConditionID"`
	PrependCID          string `gorm:"size:50" json:"PrependCID"`
	CallRecording       string `gorm:"size:20;default:'off'" json:"CallRecording"` // "off", "always", "inbound"
	MaxRingTimeSeconds  int    `gorm:"default:30" json:"MaxRingTimeSeconds"`
}

type OutboundRoute struct {
	Base
	TenantRef
	Name              string    `gorm:"size:100" json:"Name"`
	DialPattern       string    `gorm:"size:100;not null" json:"DialPattern"`
	PrefixStrip       int       `gorm:"default:0" json:"PrefixStrip"`
	PrependDigits     string    `gorm:"size:50" json:"PrependDigits"`
	TrunkID           uuid.UUID `gorm:"type:uuid;not null;index" json:"TrunkID"`
	FallbackTrunkID   string    `gorm:"size:100" json:"FallbackTrunkID"`
	Priority          int       `gorm:"default:1" json:"Priority"`
	OutboundCIDNumber string    `gorm:"size:50" json:"OutboundCIDNumber"`
	OutboundCIDName   string    `gorm:"size:100" json:"OutboundCIDName"`
	EmergencyRoute    bool      `gorm:"default:false" json:"EmergencyRoute"`
}

type Trunk struct {
	Base
	TenantRef
	ProviderName   string `gorm:"size:100;not null" json:"ProviderName"`
	SIPServer      string `gorm:"size:255;not null" json:"SIPServer"`
	SIPPort        int    `gorm:"default:5060" json:"SIPPort"`
	Transport      string `gorm:"size:20;default:'UDP'" json:"Transport"` // UDP, TCP, TLS
	AuthMethod     string `gorm:"size:50;not null" json:"AuthMethod"` // IP_ACL or USERPASS
	Username       string `gorm:"size:100" json:"Username,omitempty"`
	Password       string `gorm:"size:100" json:"Password,omitempty"`
	Register       bool   `gorm:"default:false" json:"Register"`
	MaxChannels    int    `gorm:"default:10" json:"MaxChannels"`
	Codecs         string `gorm:"size:100;default:'PCMU,PCMA,G722'" json:"Codecs"`
	FromDomain     string `gorm:"size:100" json:"FromDomain"`
	OutboundProxy  string `gorm:"size:100" json:"OutboundProxy"`
	CallerIDNumber string `gorm:"size:50" json:"CallerIDNumber"`
	IsActive       bool   `gorm:"default:true" json:"IsActive"`
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
	Name               string      `gorm:"size:100;not null" json:"Name"`
	Extension          string      `gorm:"size:50;not null;uniqueIndex:idx_tenant_ivr_ext" json:"Extension"`
	Greeting           string      `gorm:"size:255;not null" json:"Greeting"` // Path to wav or text for TTS
	Timeout            int         `gorm:"default:10000" json:"Timeout"`     // Milliseconds
	MaxFailures        int         `gorm:"default:3" json:"MaxFailures"`
	TimeoutDestType    string      `gorm:"size:50" json:"TimeoutDestType"`
	TimeoutDestTarget  string      `gorm:"size:100" json:"TimeoutDestTarget"`
	InvalidDestType    string      `gorm:"size:50" json:"InvalidDestType"`
	InvalidDestTarget  string      `gorm:"size:100" json:"InvalidDestTarget"`
	AllowDirectDial    bool        `gorm:"default:true" json:"AllowDirectDial"`
	PinCode            string      `gorm:"size:50" json:"PinCode"`
	SiteID             *uuid.UUID  `gorm:"type:uuid" json:"SiteID"`
	AnnounceOnly       bool        `gorm:"default:false" json:"AnnounceOnly"`
	TimeoutDestination string      `gorm:"size:150;default:'End Call'" json:"TimeoutDestination"`
	SoundFile          string      `gorm:"size:255" json:"SoundFile"`
	Choices            []IVRChoice `gorm:"foreignKey:MenuID" json:"Choices"`
}

// IVRChoice defines digit mapping for a menu
type IVRChoice struct {
	Base
	MenuID      uuid.UUID `gorm:"type:uuid;not null;index" json:"MenuID"`
	Digits      string    `gorm:"size:10;not null" json:"Digits"`
	Name        string    `gorm:"size:100" json:"Name"`
	Action      string    `gorm:"size:50;not null" json:"Action"`
	Destination string    `gorm:"size:255;not null" json:"Destination"`
}

// CallQueue defines ACD queue matching enterprise portal layout
type CallQueue struct {
	Base
	TenantRef
	Name                       string        `gorm:"size:100;not null" json:"Name"`
	Extension                  string        `gorm:"size:50;not null;uniqueIndex:idx_tenant_queue_ext" json:"Extension"`
	MaxCapacity                int           `gorm:"default:10" json:"MaxCapacity"`
	RingStrategy               string        `gorm:"size:50;default:'Ring All'" json:"RingStrategy"` // Ring All, Round Robin, Least Recent
	MemberRingTimeout          int           `gorm:"default:20" json:"MemberRingTimeout"`
	MemberRetryDelay           int           `gorm:"default:10" json:"MemberRetryDelay"`
	MemberWrapUpDelay          int           `gorm:"default:0" json:"MemberWrapUpDelay"`
	InternalNumber             string        `gorm:"size:50" json:"InternalNumber"`
	PlayMusicWhileWaiting      bool          `gorm:"default:true" json:"PlayMusicWhileWaiting"`
	Music                      string        `gorm:"size:100;default:'Default'" json:"Music"`
	CallerNamePrefix           string        `gorm:"size:100" json:"CallerNamePrefix"`
	RingCadence                string        `gorm:"size:50;default:'Default'" json:"RingCadence"`
	AnnounceWaitTime           bool          `gorm:"default:false" json:"AnnounceWaitTime"`
	AnnouncePosition           bool          `gorm:"default:true" json:"AnnouncePosition"`
	AnnouncementInterval       int           `gorm:"default:60" json:"AnnouncementInterval"`
	NoAnswerTimeout            int           `gorm:"default:20" json:"NoAnswerTimeout"`
	NoAnswerDestination        string        `gorm:"size:150" json:"NoAnswerDestination"`
	ParkEnabled                bool          `gorm:"default:false" json:"ParkEnabled"`
	ParkTelephoneNumber        string        `gorm:"size:50" json:"ParkTelephoneNumber"`
	ParkResumeWhenAtPosition   int           `gorm:"default:2" json:"ParkResumeWhenAtPosition"`
	ParkAllowedMinPosition     int           `gorm:"default:2" json:"ParkAllowedMinPosition"`
	ParkResumeRetryIntervalMin int           `gorm:"default:1" json:"ParkResumeRetryIntervalMin"`
	ParkResumeMaxAttempts      int           `gorm:"default:3" json:"ParkResumeMaxAttempts"`
	ParkResumeMenuTimeoutSec   int           `gorm:"default:15" json:"ParkResumeMenuTimeoutSec"`
	ParkResumeDialTimeoutSec   int           `gorm:"default:20" json:"ParkResumeDialTimeoutSec"`
	ParkConfirmation           string        `gorm:"size:100" json:"ParkConfirmation"`
	Members                    []QueueMember `gorm:"foreignKey:QueueID" json:"Members"`
}

// QueueMember defines an assigned agent for a CallQueue
type QueueMember struct {
	Base
	QueueID   uuid.UUID `gorm:"type:uuid;not null;index" json:"QueueID"`
	Member    string    `gorm:"size:100;not null" json:"Member"` // Display Name / Agent Name
	Extension string    `gorm:"size:50" json:"Extension"`
	Penalty   int       `gorm:"default:0" json:"Penalty"`
	Disabled  bool      `gorm:"default:false" json:"Disabled"`
}

// CallDistributionGroup defines a Ring Group / Hunt Group matching enterprise portal
type CallDistributionGroup struct {
	Base
	TenantRef
	Name                string        `gorm:"size:100;not null" json:"Name"`
	SiteID              *uuid.UUID    `gorm:"type:uuid" json:"SiteID"`
	SiteName            string        `gorm:"size:100" json:"SiteName"`
	InternalNumber      string        `gorm:"size:50;not null;uniqueIndex:idx_tenant_group_ext" json:"InternalNumber"`
	RingCadence         string        `gorm:"size:50;default:'Default'" json:"RingCadence"`
	CallerNamePrefix    string        `gorm:"size:100" json:"CallerNamePrefix"`
	NoAnswerTimeout     int           `gorm:"default:15" json:"NoAnswerTimeout"`
	NoAnswerDestination string        `gorm:"size:150" json:"NoAnswerDestination"`
	Members             []GroupMember `gorm:"foreignKey:GroupID" json:"Members"`
}

// GroupMember defines an extension in a CallDistributionGroup with ringing delay
type GroupMember struct {
	Base
	GroupID   uuid.UUID `gorm:"type:uuid;not null;index" json:"GroupID"`
	Member    string    `gorm:"size:100;not null" json:"Member"` // Display Name / Extension
	Extension string    `gorm:"size:50" json:"Extension"`
	DelaySec  int       `gorm:"default:0" json:"DelaySec"`
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

// ButtonProfile defines programmable DSS/BLF line keys for desk phones
type ButtonProfile struct {
	Base
	TenantRef
	Name        string `gorm:"size:100;not null" json:"Name"`
	DeviceModel string `gorm:"size:100;default:'Generic SIP Handset'" json:"DeviceModel"`
	Description string `gorm:"size:255" json:"Description"`
	Keys        string `gorm:"type:text" json:"Keys"` // JSON string of [{key: 1, type: "BLF", value: "1001", label: "Front Desk"}]
}

// Contact represents an entry in the centralized company phonebook/directory
type Contact struct {
	Base
	TenantRef
	FirstName    string `gorm:"size:100;not null" json:"FirstName"`
	LastName     string `gorm:"size:100" json:"LastName"`
	Company      string `gorm:"size:150" json:"Company"`
	PhoneNumber  string `gorm:"size:50;not null" json:"PhoneNumber"`
	MobileNumber string `gorm:"size:50" json:"MobileNumber"`
	Email        string `gorm:"size:150" json:"Email"`
	SpeedDial    string `gorm:"size:10" json:"SpeedDial"`
	Group        string `gorm:"size:50;default:'General'" json:"Group"`
}

// Site represents a company branch office / physical location
type Site struct {
	Base
	TenantRef
	Name             string `gorm:"size:100;not null" json:"Name"`
	SiteCode         string `gorm:"size:10;not null" json:"SiteCode"` // e.g. 71, 72
	Timezone         string `gorm:"size:50;default:'UTC'" json:"Timezone"`
	CallerIDNumber   string `gorm:"size:50" json:"CallerIDNumber"`
	EmergencyAddress string `gorm:"size:255" json:"EmergencyAddress"`
	IsDefault        bool   `gorm:"default:false" json:"IsDefault"`
}

// TelephoneNumber represents a phone number / DID assigned to the tenant
type TelephoneNumber struct {
	Base
	TenantRef
	Number            string     `gorm:"size:50;not null;uniqueIndex:idx_tenant_telnum" json:"Number"`
	Description       string     `gorm:"size:100" json:"Description"`
	DestinationType   string     `gorm:"size:50;default:'extension'" json:"DestinationType"` // extension, ivr, group, queue, voicemail
	DestinationTarget string     `gorm:"size:100" json:"DestinationTarget"`
	OutboundCLI       bool       `gorm:"default:true" json:"OutboundCLI"`
	SiteID            *uuid.UUID `gorm:"type:uuid" json:"SiteID"`
	Provider          string     `gorm:"size:100;default:'Default Carrier'" json:"Provider"`
}

// CallDivert manages call forwarding policies per extension
type CallDivert struct {
	Base
	TenantRef
	ExtensionNumber string `gorm:"size:50;not null;uniqueIndex:idx_tenant_divert_ext" json:"ExtensionNumber"`
	AlwaysEnabled   bool   `gorm:"default:false" json:"AlwaysEnabled"`
	AlwaysTarget    string `gorm:"size:50" json:"AlwaysTarget"`
	BusyEnabled     bool   `gorm:"default:false" json:"BusyEnabled"`
	BusyTarget      string `gorm:"size:50" json:"BusyTarget"`
	NoAnswerEnabled bool   `gorm:"default:false" json:"NoAnswerEnabled"`
	NoAnswerTarget  string `gorm:"size:50" json:"NoAnswerTarget"`
	NoAnswerTimeout int    `gorm:"default:20" json:"NoAnswerTimeout"`
	OfflineEnabled  bool   `gorm:"default:false" json:"OfflineEnabled"`
	OfflineTarget   string `gorm:"size:50" json:"OfflineTarget"`
}

// Device manages hardware phones, MAC addresses, and provisioning status
type Device struct {
	Base
	TenantRef
	MACAddress      string     `gorm:"size:50;not null;uniqueIndex:idx_tenant_dev_mac" json:"MACAddress"`
	BrandModel      string     `gorm:"size:100;not null" json:"BrandModel"`
	ExtensionNumber string     `gorm:"size:50;not null" json:"ExtensionNumber"`
	IPAddress       string     `gorm:"size:50" json:"IPAddress"`
	FirmwareVersion string     `gorm:"size:50" json:"FirmwareVersion"`
	Status          string     `gorm:"size:50;default:'online'" json:"Status"` // online, offline, provisioning
	ProfileID       *uuid.UUID `gorm:"type:uuid" json:"ProfileID"`
}
