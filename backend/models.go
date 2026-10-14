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
	Summary        string    `gorm:"type:text" json:"Summary"`
	Sentiment      string    `gorm:"size:50;default:'neutral'" json:"Sentiment"`
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
	Extension     string `gorm:"size:50;not null" json:"Extension"`
	CallerID      string `gorm:"size:100" json:"CallerID"`
	Duration      int    `gorm:"default:0" json:"Duration"`
	FilePath      string `gorm:"size:255;not null" json:"FilePath"`
	IsRead        bool   `gorm:"default:false" json:"IsRead"`
	Transcription string `gorm:"type:text" json:"Transcription"`
	Summary       string `gorm:"type:text" json:"Summary"`
	Sentiment     string `gorm:"size:50;default:'neutral'" json:"Sentiment"`
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
	Name              string `gorm:"size:100;not null" json:"name"`
	DaysOfWeek        string `gorm:"size:100;default:'mon-fri'" json:"days_of_week"`
	TimeStart         string `gorm:"size:10;default:'09:00'" json:"time_start"`
	TimeEnd           string `gorm:"size:10;default:'17:00'" json:"time_end"`
	MatchedDestType   string `gorm:"size:50;default:'ivr'" json:"matched_dest_type"`
	MatchedDestVal    string `gorm:"size:100;not null" json:"matched_dest_val"`
	UnmatchedDestType string `gorm:"size:50;default:'voicemail'" json:"unmatched_dest_type"`
	UnmatchedDestVal  string `gorm:"size:100;not null" json:"unmatched_dest_val"`
	IsActive          bool   `gorm:"default:true" json:"is_active"`
}

// Holiday defines scheduled holidays
type Holiday struct {
	Base
	TenantRef
	Name      string `gorm:"size:100;not null" json:"name"`
	StartDate string `gorm:"size:50;not null" json:"start_date"`
	EndDate   string `gorm:"size:50;not null" json:"end_date"`
	DestType  string `gorm:"size:50;default:'voicemail'" json:"dest_type"`
	DestVal   string `gorm:"size:100;not null" json:"dest_val"`
	IsActive  bool   `gorm:"default:true" json:"is_active"`
}

// OutboundRestriction blocks dialing specific prefixes
type OutboundRestriction struct {
	Base
	TenantRef
	Name        string `gorm:"size:100;not null" json:"name"`
	Pattern     string `gorm:"size:100;not null" json:"pattern"` // e.g. 00X., 090X.
	BarredTypes string `gorm:"size:50;default:'international'" json:"barred_types"`
	AllowedExts string `gorm:"size:255" json:"allowed_exts"`
	PinCode     string `gorm:"size:50" json:"pin_code"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
}

// AutoCLIP routes returning callers directly to the extension that last phoned them
type AutoCLIP struct {
	Base
	TenantRef
	CallerNumber      string    `gorm:"size:50;not null;index" json:"caller_number"`
	OriginalCalleeExt string    `gorm:"size:50;not null" json:"original_callee_ext"`
	TrunkName         string    `gorm:"size:100;default:'Default SIP Trunk'" json:"trunk_name"`
	ExpiresAt         time.Time `json:"expires_at"`
}

// PagingGroup defines intercom/paging broadcast groups
type PagingGroup struct {
	Base
	TenantRef
	Name       string `gorm:"size:100;not null" json:"name"`
	Number     string `gorm:"size:50;not null;uniqueIndex:idx_tenant_paging_num" json:"number"`
	DuplexMode string `gorm:"size:20;default:'one-way'" json:"duplex_mode"` // "one-way" or "two-way"
	Members    string `gorm:"type:text;not null" json:"members"`
}

// PickupGroup allows users to grab ringing calls with *8
type PickupGroup struct {
	Base
	TenantRef
	Name    string `gorm:"size:100;not null" json:"name"`
	Code    string `gorm:"size:50;not null" json:"code"` // e.g. *8
	Members string `gorm:"type:text;not null" json:"members"`
}

// AudioPrompt manages MoH playlists, system prompts, and custom audio files
type AudioPrompt struct {
	Base
	TenantRef
	Name        string `gorm:"size:100;not null" json:"name"`
	Type        string `gorm:"size:50;default:'custom'" json:"type"` // "moh", "custom", "system"
	Description string `gorm:"size:255" json:"description"`
	FilePath    string `gorm:"size:255;not null" json:"file_path"`
	IsDefault   bool   `gorm:"default:false" json:"is_default"`
}

// FeatureCode lists and configures dialed feature codes
type FeatureCode struct {
	Base
	TenantRef
	Name        string `gorm:"size:100;not null" json:"name"`
	Code        string `gorm:"size:20;not null" json:"code"`
	DefaultCode string `gorm:"size:20;not null" json:"default_code"`
	Description string `gorm:"size:255" json:"description"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
}

// RecordingRule configures automated call recording
type RecordingRule struct {
	Base
	TenantRef
	Name        string `gorm:"size:100;not null" json:"name"`
	TargetType  string `gorm:"size:50;default:'all'" json:"target_type"` // "all", "extension", "queue", "trunk"
	TargetID    string `gorm:"size:100" json:"target_id"`
	Direction   string `gorm:"size:50;default:'both'" json:"direction"` // "inbound", "outbound", "both"
	StorageDays int    `gorm:"default:90" json:"storage_days"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
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

// SIPSettings holds comprehensive PBX, SIP signaling, RTP media, NAT/ICE, timers, encryption, and QoS configuration
type SIPSettings struct {
	Base
	TenantRef

	// 1. Signaling & Listening Ports
	UserAgent           string `gorm:"size:100;default:'MSIP-PBX-Core/2.0'" json:"user_agent"`
	SIPPort             int    `gorm:"default:5060" json:"sip_port"`
	SIPTLSPort          int    `gorm:"default:5061" json:"sip_tls_port"`
	ExternalSIPPort     int    `gorm:"default:5080" json:"external_sip_port"`
	ExternalSIPTLSPort  int    `gorm:"default:5081" json:"external_sip_tls_port"`
	WSSPort             int    `gorm:"default:7443" json:"wss_port"`
	TCPKeepAliveSec     int    `gorm:"default:30" json:"tcp_keepalive_sec"`
	CompactHeaders      bool   `gorm:"default:false" json:"compact_headers"`
	AllowContactRewrite bool   `gorm:"default:true" json:"allow_contact_rewrite"`

	// 2. RTP Media & Jitter Buffer
	RTPPortStart          int    `gorm:"default:16384" json:"rtp_port_start"`
	RTPPortEnd            int    `gorm:"default:32768" json:"rtp_port_end"`
	RTPTimeoutSec         int    `gorm:"default:30" json:"rtp_timeout_sec"`
	RTPHoldTimeoutSec     int    `gorm:"default:1800" json:"rtp_hold_timeout_sec"`
	JitterBufferMode      string `gorm:"size:20;default:'adaptive'" json:"jitter_buffer_mode"`
	JitterBufferInitialMs int    `gorm:"default:60" json:"jitter_buffer_initial_ms"`
	JitterBufferMaxMs     int    `gorm:"default:200" json:"jitter_buffer_max_ms"`
	DTMFMode              string `gorm:"size:20;default:'rfc2833'" json:"dtmf_mode"`
	DTMFDurationMs        int    `gorm:"default:100" json:"dtmf_duration_ms"`
	AudioCodecs           string `gorm:"size:200;default:'OPUS,G722,PCMA,PCMU'" json:"audio_codecs"`
	VideoCodecs           string `gorm:"size:100;default:'VP8,H264'" json:"video_codecs"`
	EnableCNG             bool   `gorm:"default:true" json:"enable_cng"`

	// 3. NAT, STUN, TURN & ICE Traversal
	NATMode            string `gorm:"size:20;default:'auto'" json:"nat_mode"`
	ExternalIPv4       string `gorm:"size:100" json:"external_ipv4"`
	ExternalIPv6       string `gorm:"size:100" json:"external_ipv6"`
	LocalNetworks      string `gorm:"size:255;default:'127.0.0.1/32, 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16'" json:"local_networks"`
	STUNServer         string `gorm:"size:255;default:'stun:stun.l.google.com:19302'" json:"stun_server"`
	TURNServer         string `gorm:"size:255" json:"turn_server"`
	TURNUser           string `gorm:"size:100" json:"turn_user"`
	TURNPassword       string `gorm:"size:100" json:"turn_password"`
	EnableICE          bool   `gorm:"default:true" json:"enable_ice"`
	AggressiveICE      bool   `gorm:"default:true" json:"aggressive_ice"`
	NATPingIntervalSec int    `gorm:"default:20" json:"nat_ping_interval_sec"`

	// 4. SIP Session Timers & Registration TTL
	SessionExpiresSec  int    `gorm:"default:1800" json:"session_expires_sec"`
	MinSESec           int    `gorm:"default:120" json:"min_se_sec"`
	SessionRefresher   string `gorm:"size:20;default:'auto'" json:"session_refresher"`
	RegisterMinSec     int    `gorm:"default:60" json:"register_min_sec"`
	RegisterMaxSec     int    `gorm:"default:3600" json:"register_max_sec"`
	RegisterDefaultSec int    `gorm:"default:3600" json:"register_default_sec"`
	RegisterRetrySec   int    `gorm:"default:30" json:"register_retry_sec"`
	MaxChannels        int    `gorm:"default:250" json:"max_channels"`

	// 5. Security, TLS & SRTP Encryption
	SRTPMode            string `gorm:"size:20;default:'optional'" json:"srtp_mode"`
	SRTPSuites          string `gorm:"size:255;default:'AEAD_AES_256_GCM,AES_CM_128_HMAC_SHA1_80,AES_CM_128_HMAC_SHA1_32'" json:"srtp_suites"`
	TLSVersionMin       string `gorm:"size:10;default:'1.2'" json:"tls_version_min"`
	TLSVerifyPeer       bool   `gorm:"default:false" json:"tls_verify_peer"`
	AuthRealm           string `gorm:"size:100;default:'msip.local'" json:"auth_realm"`
	ChallengeAllInvites bool   `gorm:"default:true" json:"challenge_all_invites"`

	// 6. Quality of Service (QoS) & DSCP Marking
	DSCPAudio     int `gorm:"default:46" json:"dscp_audio"`
	DSCPVideo     int `gorm:"default:34" json:"dscp_video"`
	DSCPSignaling int `gorm:"default:26" json:"dscp_signaling"`
}

// DefaultSIPSettings returns a pre-configured SIPSettings instance with industry best-practice values
func DefaultSIPSettings(tenantID uuid.UUID) SIPSettings {
	return SIPSettings{
		Base:                Base{ID: uuid.New()},
		TenantRef:           TenantRef{TenantID: tenantID},
		UserAgent:           "MSIP-PBX-Core/2.0",
		SIPPort:             5060,
		SIPTLSPort:          5061,
		ExternalSIPPort:     5080,
		ExternalSIPTLSPort:  5081,
		WSSPort:             7443,
		TCPKeepAliveSec:     30,
		CompactHeaders:      false,
		AllowContactRewrite: true,

		RTPPortStart:          16384,
		RTPPortEnd:            32768,
		RTPTimeoutSec:         30,
		RTPHoldTimeoutSec:     1800,
		JitterBufferMode:      "adaptive",
		JitterBufferInitialMs: 60,
		JitterBufferMaxMs:     200,
		DTMFMode:              "rfc2833",
		DTMFDurationMs:        100,
		AudioCodecs:           "OPUS,G722,PCMA,PCMU",
		VideoCodecs:           "VP8,H264",
		EnableCNG:             true,

		NATMode:            "auto",
		ExternalIPv4:       "",
		ExternalIPv6:       "",
		LocalNetworks:      "127.0.0.1/32, 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16",
		STUNServer:         "stun:stun.l.google.com:19302",
		TURNServer:         "",
		TURNUser:           "",
		TURNPassword:       "",
		EnableICE:          true,
		AggressiveICE:      true,
		NATPingIntervalSec: 20,

		SessionExpiresSec:  1800,
		MinSESec:           120,
		SessionRefresher:   "auto",
		RegisterMinSec:     60,
		RegisterMaxSec:     3600,
		RegisterDefaultSec: 3600,
		RegisterRetrySec:   30,
		MaxChannels:        250,

		SRTPMode:            "optional",
		SRTPSuites:          "AEAD_AES_256_GCM,AES_CM_128_HMAC_SHA1_80,AES_CM_128_HMAC_SHA1_32",
		TLSVersionMin:       "1.2",
		TLSVerifyPeer:       false,
		AuthRealm:           "msip.local",
		ChallengeAllInvites: true,

		DSCPAudio:     46,
		DSCPVideo:     34,
		DSCPSignaling: 26,
	}
}

// AIAssistant defines a conversational AI telephone agent for virtual reception, booking, and natural language IVR
type AIAssistant struct {
	Base
	TenantRef
	Name               string `gorm:"size:100;not null" json:"name"`
	Extension          string `gorm:"size:50;not null;uniqueIndex:idx_tenant_ai_ext" json:"extension"`
	ExtensionNumber    string `gorm:"size:50" json:"extension_number,omitempty"`
	RoleType           string `gorm:"size:50;default:'receptionist'" json:"role_type"` // receptionist, booking, ivr, support
	RoleDescription    string `gorm:"size:255" json:"role_description,omitempty"`
	Greeting           string `gorm:"size:500" json:"greeting"`
	GreetingMessage    string `gorm:"size:500" json:"greeting_message,omitempty"`
	SystemPrompt       string `gorm:"type:text" json:"system_prompt"`
	LLMProvider        string `gorm:"size:50;default:'openai'" json:"llm_provider"` // openai, claude, gemini, ollama
	LLMModel           string `gorm:"size:100;default:'gpt-4o'" json:"llm_model"`
	VoiceProvider      string `gorm:"size:50;default:'openai'" json:"voice_provider"` // openai, elevenlabs, deepgram, system
	TTSProvider        string `gorm:"size:50" json:"tts_provider,omitempty"`
	VoiceName          string `gorm:"size:50;default:'alloy'" json:"voice_name"`
	VoiceID            string `gorm:"size:50" json:"voice_id,omitempty"`
	Language           string `gorm:"size:20;default:'en-US'" json:"language"`
	EnableBooking      bool   `gorm:"default:true" json:"enable_booking"`
	EnableTransfer     bool   `gorm:"default:true" json:"enable_transfer"`
	TransferTarget     string `gorm:"size:50;default:'1001'" json:"transfer_target"`
	EnableVoicemail    bool   `gorm:"default:true" json:"enable_voicemail"`
	MaxCallDurationSec int    `gorm:"default:600" json:"max_call_duration_sec"`
	IsActive           bool   `gorm:"default:true" json:"is_active"`
}

func (a *AIAssistant) EffectiveGreeting() string {
	if a.Greeting != "" {
		return a.Greeting
	}
	if a.GreetingMessage != "" {
		return a.GreetingMessage
	}
	return "Hello! How may I assist you today?"
}

func (a *AIAssistant) EffectiveVoice() string {
	if a.VoiceName != "" {
		return a.VoiceName
	}
	if a.VoiceID != "" {
		return a.VoiceID
	}
	return "alloy"
}

func (a *AIAssistant) EffectiveVoiceProvider() string {
	if a.VoiceProvider != "" {
		return a.VoiceProvider
	}
	if a.TTSProvider != "" {
		return a.TTSProvider
	}
	return "openai"
}

func (a *AIAssistant) EffectiveExtension() string {
	if a.Extension != "" {
		return a.Extension
	}
	return a.ExtensionNumber
}

// AIBooking manages appointments scheduled by the conversational AI agent
type AIBooking struct {
	Base
	TenantRef
	AssistantID     uuid.UUID `gorm:"type:uuid;index" json:"assistant_id"`
	AssistantName   string    `gorm:"size:100" json:"assistant_name"`
	CallerNumber    string    `gorm:"size:50;not null" json:"caller_number"`
	CallerName      string    `gorm:"size:100;not null" json:"caller_name"`
	ServiceType     string    `gorm:"size:100;default:'General Appointment'" json:"service_type"`
	BookingDate     string    `gorm:"size:20;not null" json:"booking_date"` // YYYY-MM-DD
	BookingTime     string    `gorm:"size:10;not null" json:"booking_time"` // HH:MM
	DurationMinutes int       `gorm:"default:30" json:"duration_minutes"`
	Status          string    `gorm:"size:50;default:'confirmed'" json:"status"` // confirmed, pending, completed, cancelled
	Notes           string    `gorm:"type:text" json:"notes"`
	CallUUID        string    `gorm:"size:100" json:"call_uuid"`
}

// AICallSession tracks transcripts, summaries, sentiment, and executed tools for AI calls
type AICallSession struct {
	Base
	TenantRef
	AssistantID   uuid.UUID `gorm:"type:uuid;index" json:"assistant_id"`
	AssistantName string    `gorm:"size:100" json:"assistant_name"`
	CallUUID      string    `gorm:"size:100;index" json:"call_uuid"`
	CallerNumber  string    `gorm:"size:50" json:"caller_number"`
	Summary       string    `gorm:"type:text" json:"summary"`
	Sentiment     string    `gorm:"size:50;default:'neutral'" json:"sentiment"` // positive, neutral, negative
	ActionTaken   string    `gorm:"size:50;default:'conversation'" json:"action_taken"` // booking_created, transferred, message_taken, info_provided
	MessagesJSON  string    `gorm:"type:text" json:"messages_json"`
	DurationSec   int       `gorm:"default:0" json:"duration_sec"`
}

// AITranscriptionSettings holds global configuration for STT transcription, TTS synthesis, and LLM backends
type AITranscriptionSettings struct {
	Base
	TenantRef
	AutoTranscribeVoicemail  bool   `gorm:"default:true" json:"auto_transcribe_voicemail"`
	AutoTranscribeRecordings bool   `gorm:"default:true" json:"auto_transcribe_recordings"`
	STTProvider              string `gorm:"size:50;default:'whisper'" json:"stt_provider"` // whisper, deepgram, google, vosk
	STTModel                 string `gorm:"size:50;default:'whisper-1'" json:"stt_model"`
	STTAPIKey                string `gorm:"size:255" json:"stt_api_key"`
	TTSProvider              string `gorm:"size:50;default:'openai'" json:"tts_provider"` // openai, elevenlabs, deepgram, system
	TTSAPIKey                string `gorm:"size:255" json:"tts_api_key"`
	TTSDefaultVoice          string `gorm:"size:50;default:'alloy'" json:"tts_default_voice"`
	LLMProvider              string `gorm:"size:50;default:'openai'" json:"llm_provider"` // openai, claude, gemini, ollama
	LLMAPIKey                string `gorm:"size:255" json:"llm_api_key"`
	LLMDefaultModel          string `gorm:"size:100;default:'gpt-4o'" json:"llm_default_model"`
	LLMBaseURL               string `gorm:"size:255" json:"llm_base_url"` // for local Ollama/vLLM http://localhost:11434/v1
	SentimentAnalysis        bool   `gorm:"default:true" json:"sentiment_analysis"`
	EmailTranscripts         bool   `gorm:"default:false" json:"email_transcripts"`
	NotificationEmail        string `gorm:"size:150" json:"notification_email"`
}

// DefaultAITranscriptionSettings returns recommended default settings for AI & transcription
func DefaultAITranscriptionSettings(tenantID uuid.UUID) AITranscriptionSettings {
	return AITranscriptionSettings{
		Base:                     Base{ID: uuid.New()},
		TenantRef:                TenantRef{TenantID: tenantID},
		AutoTranscribeVoicemail:  true,
		AutoTranscribeRecordings: true,
		STTProvider:              "whisper",
		STTModel:                 "whisper-1",
		STTAPIKey:                "",
		TTSProvider:              "openai",
		TTSAPIKey:                "",
		TTSDefaultVoice:          "alloy",
		LLMProvider:              "openai",
		LLMAPIKey:                "",
		LLMDefaultModel:          "gpt-4o",
		LLMBaseURL:               "https://api.openai.com/v1",
		SentimentAnalysis:        true,
		EmailTranscripts:         false,
		NotificationEmail:        "",
	}
}
