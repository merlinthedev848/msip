package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB initializes database connection with support for Postgres and SQLite.
func InitDB() {
	dsn := os.Getenv("DATABASE_URL")
	driver := strings.ToLower(os.Getenv("DATABASE_DRIVER"))

	var dialector gorm.Dialector

	if driver == "sqlite" || strings.HasPrefix(dsn, "sqlite://") || strings.Contains(dsn, ".db") || strings.Contains(dsn, ":memory:") {
		sqlitePath := dsn
		if strings.HasPrefix(sqlitePath, "sqlite://") {
			sqlitePath = strings.TrimPrefix(sqlitePath, "sqlite://")
		}
		if sqlitePath == "" {
			sqlitePath = "msip.db"
		}
		log.Printf("Connecting to SQLite database: %s", sqlitePath)
		dialector = sqlite.Open(sqlitePath)
	} else {
		if dsn == "" {
			dsn = "host=localhost user=msip_user password=msip_secure dbname=msip_pbx port=5432 sslmode=disable TimeZone=UTC"
		}
		log.Printf("Connecting to PostgreSQL database")
		dialector = postgres.Open(dsn)
	}

	var err error
	DB, err = gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := DB.DB()
	if err == nil {
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	// In Postgres, ensure pgcrypto is enabled
	if dialector.Name() == "postgres" {
		DB.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto;")
	}

	log.Println("Database connection established. Running auto-migrations...")
	err = DB.AutoMigrate(
		&Tenant{},
		&User{},
		&Extension{},
		&InboundRoute{},
		&OutboundRoute{},
		&Trunk{},
		&CDR{},
		&IVRMenu{},
		&IVRChoice{},
		&Voicemail{},
		&VideoRoom{},
		&IPCamera{},
		&TeamChat{},
		&SMSMessage{},
		&APICredential{},
		&Webhook{},
		&ThrottleRule{},
		&FirewallRule{},
		&FraudRule{},
		&TimeCondition{},
		&Holiday{},
		&OutboundRestriction{},
		&AutoCLIP{},
		&PagingGroup{},
		&PickupGroup{},
		&AudioPrompt{},
		&FeatureCode{},
		&RecordingRule{},
		&ButtonProfile{},
		&Contact{},
		&Site{},
		&TelephoneNumber{},
		&CallDivert{},
		&Device{},
		&CallQueue{},
		&QueueMember{},
		&CallDistributionGroup{},
		&GroupMember{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}

	SeedDefaults()
}

// SeedDefaults populates initial organization, admin user, test extension, and feature codes if empty.
func SeedDefaults() {
	var count int64
	DB.Model(&Tenant{}).Count(&count)
	if count == 0 {
		defaultTenant := Tenant{
			Name:   "Default Organization",
			Domain: "pbx.local",
		}
		if err := DB.Create(&defaultTenant).Error; err != nil {
			log.Printf("Error seeding tenant: %v", err)
			return
		}
		log.Printf("Seeded default tenant: %s (%s)", defaultTenant.Name, defaultTenant.ID.String())

		hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		defaultUser := User{
			TenantID:     defaultTenant.ID,
			Email:        "admin@pbx.local",
			PasswordHash: string(hash),
			Role:         "superadmin",
		}
		if err := DB.Create(&defaultUser).Error; err != nil {
			log.Printf("Error seeding default user: %v", err)
		} else {
			log.Printf("Seeded default admin user: %s", defaultUser.Email)
		}

		defaultExt := Extension{
			TenantRef:        TenantRef{TenantID: defaultTenant.ID},
			ExtensionNumber:  "1001",
			PasswordHash:     "testpass",
			IsActive:         true,
			ProvisioningCode: "10011001",
		}
		if err := DB.Create(&defaultExt).Error; err != nil {
			log.Printf("Error seeding default extension: %v", err)
		} else {
			log.Printf("Seeded default extension: %s", defaultExt.ExtensionNumber)
		}

		// Seed standard telephony feature codes
		featureCodes := []FeatureCode{
			{TenantRef: TenantRef{TenantID: defaultTenant.ID}, Name: "Check Personal Voicemail", Code: "*97", DefaultCode: "*97", Description: "Access caller's personal voicemail mailbox", IsActive: true},
			{TenantRef: TenantRef{TenantID: defaultTenant.ID}, Name: "Check General Voicemail", Code: "*98", DefaultCode: "*98", Description: "Access any voicemail box with PIN", IsActive: true},
			{TenantRef: TenantRef{TenantID: defaultTenant.ID}, Name: "Call Pickup (Any Ringing Handset)", Code: "*8", DefaultCode: "*8", Description: "Grab an incoming call ringing on another phone", IsActive: true},
			{TenantRef: TenantRef{TenantID: defaultTenant.ID}, Name: "Call Forwarding Always (Enable)", Code: "*72", DefaultCode: "*72", Description: "Forward all incoming calls unconditionally", IsActive: true},
			{TenantRef: TenantRef{TenantID: defaultTenant.ID}, Name: "Call Forwarding Always (Cancel)", Code: "*73", DefaultCode: "*73", Description: "Disable unconditional call forwarding", IsActive: true},
			{TenantRef: TenantRef{TenantID: defaultTenant.ID}, Name: "Two-Way Intercom Dialing", Code: "*0", DefaultCode: "*0", Description: "Auto-answer handsfree intercom broadcast", IsActive: true},
			{TenantRef: TenantRef{TenantID: defaultTenant.ID}, Name: "Echo Audio Latency Test", Code: "9999", DefaultCode: "9999", Description: "Two-way RTP audio loopback test", IsActive: true},
			{TenantRef: TenantRef{TenantID: defaultTenant.ID}, Name: "Video / Audio Conference Room", Code: "9000", DefaultCode: "9000", Description: "Main multi-party MCU conference bridge", IsActive: true},
			{TenantRef: TenantRef{TenantID: defaultTenant.ID}, Name: "AI Voice Dictation & Transcript", Code: "*88", DefaultCode: "*88", Description: "Record speech and transcribe with AI", IsActive: true},
			{TenantRef: TenantRef{TenantID: defaultTenant.ID}, Name: "Conversational AI Receptionist", Code: "*99", DefaultCode: "*99", Description: "Interactive conversational AI bot stream", IsActive: true},
		}
		for _, fc := range featureCodes {
			_ = DB.Create(&fc)
		}

		// Seed default business hours time condition
		defaultTC := TimeCondition{
			TenantRef:           TenantRef{TenantID: defaultTenant.ID},
			Name:                "Business Office Hours",
			DaysOfWeek:          "Mon,Tue,Wed,Thu,Fri",
			StartTime:           "09:00",
			EndTime:             "17:00",
			InOfficeDestType:    "EXTENSION",
			InOfficeTarget:      "1001",
			OutOfOfficeDestType: "VOICEMAIL",
			OutOfOfficeTarget:   "1001",
		}
		_ = DB.Create(&defaultTC)

		// Seed default Music on Hold prompt
		defaultMoH := AudioPrompt{
			TenantRef:   TenantRef{TenantID: defaultTenant.ID},
			Name:        "Default Telecom Harmony",
			Category:    "MOH",
			FileName:    "local_stream://moh",
			Description: "Default FreeSWITCH high-fidelity stereo acoustic hold music",
		}
		_ = DB.Create(&defaultMoH)

		// Seed default Headquarters Site
		defaultSite := Site{
			TenantRef:        TenantRef{TenantID: defaultTenant.ID},
			Name:             "Main Corporate HQ",
			SiteCode:         "10",
			Timezone:         "America/New_York",
			CallerIDNumber:   "+18005550100",
			EmergencyAddress: "100 Telecom Way, Suite 400, New York, NY 10001",
			IsDefault:        true,
		}
		_ = DB.Create(&defaultSite)

		// Seed sample company directory contact
		defaultContact := Contact{
			TenantRef:    TenantRef{TenantID: defaultTenant.ID},
			FirstName:    "Telecom",
			LastName:     "Support",
			Company:      "Internal Operations",
			PhoneNumber:  "1001",
			MobileNumber: "+15550199123",
			Email:        "support@pbx.local",
			SpeedDial:    "01",
			Group:        "Technical",
		}
		_ = DB.Create(&defaultContact)

		// Seed default button profile template
		defaultBP := ButtonProfile{
			TenantRef:   TenantRef{TenantID: defaultTenant.ID},
			Name:        "Executive IP Phone (10 Keys)",
			DeviceModel: "Standard SIP Phone",
			Description: "Default DSS/BLF layout with front desk line keys and park shortcuts",
			Keys:        `[{"key":1,"type":"Line","value":"1","label":"Line 1"},{"key":2,"type":"BLF","value":"1001","label":"Operator"},{"key":3,"type":"SpeedDial","value":"*97","label":"Voicemail"},{"key":4,"type":"SpeedDial","value":"*8","label":"Pickup"}]`,
		}
		_ = DB.Create(&defaultBP)

		// Seed primary Telephone Number (DID)
		defaultTel := TelephoneNumber{
			TenantRef:         TenantRef{TenantID: defaultTenant.ID},
			Number:            "+18005550100",
			Description:       "Primary Inbound Toll-Free Number",
			DestinationType:   "extension",
			DestinationTarget: "1001",
			OutboundCLI:       true,
			Provider:          "Primary Carrier Trunk",
		}
		_ = DB.Create(&defaultTel)

		// Seed sample Device
		defaultDev := Device{
			TenantRef:       TenantRef{TenantID: defaultTenant.ID},
			MACAddress:      "00:15:65:8A:2B:11",
			BrandModel:      "Standard HD SIP Phone",
			ExtensionNumber: "1001",
			IPAddress:       "192.168.1.105",
			FirmwareVersion: "v2.8.4",
			Status:          "online",
		}
		_ = DB.Create(&defaultDev)

		// Seed default call divert for extension 1001
		defaultDivert := CallDivert{
			TenantRef:       TenantRef{TenantID: defaultTenant.ID},
			ExtensionNumber: "1001",
			NoAnswerEnabled: true,
			NoAnswerTarget:  "1001",
			NoAnswerTimeout: 20,
		}
		_ = DB.Create(&defaultDivert)

		// Seed sample Call Distribution Group matching Screenshot 4
		defaultGroup := CallDistributionGroup{
			TenantRef:           TenantRef{TenantID: defaultTenant.ID},
			Name:                "Front door newcastle",
			SiteName:            "Agilico Systems",
			InternalNumber:      "990",
			RingCadence:         "Default",
			CallerNamePrefix:    "Front Door (9)",
			NoAnswerTimeout:     15,
			NoAnswerDestination: "Front Door Failover newcastle",
		}
		if err := DB.Create(&defaultGroup).Error; err == nil {
			sampleMembers := []GroupMember{
				{GroupID: defaultGroup.ID, Member: "Nicola Simpson", Extension: "1001", DelaySec: 0},
				{GroupID: defaultGroup.ID, Member: "Anne Mitchell Hall", Extension: "1002", DelaySec: 0},
				{GroupID: defaultGroup.ID, Member: "Michelle Boomer", Extension: "1003", DelaySec: 0},
				{GroupID: defaultGroup.ID, Member: "spare7985", Extension: "7985", DelaySec: 0},
				{GroupID: defaultGroup.ID, Member: "Michela Walsh", Extension: "1004", DelaySec: 0},
				{GroupID: defaultGroup.ID, Member: "Caitlin Redhead", Extension: "1005", DelaySec: 0},
			}
			for _, m := range sampleMembers {
				_ = DB.Create(&m)
			}
		}

		// Seed sample Call Queue matching Screenshot 3
		defaultQueue := CallQueue{
			TenantRef:                  TenantRef{TenantID: defaultTenant.ID},
			Name:                       "CDS - Accounts",
			Extension:                  "8001",
			MaxCapacity:                10,
			RingStrategy:               "Ring All",
			MemberRingTimeout:          20,
			MemberRetryDelay:           10,
			MemberWrapUpDelay:          0,
			InternalNumber:             "991",
			PlayMusicWhileWaiting:      true,
			Music:                      "Default",
			CallerNamePrefix:           "Accounts",
			RingCadence:                "Default",
			AnnounceWaitTime:           false,
			AnnouncePosition:           true,
			AnnouncementInterval:       60,
			NoAnswerTimeout:            20,
			NoAnswerDestination:        "Agilico Accounts Menu",
			ParkEnabled:                false,
			ParkResumeWhenAtPosition:   2,
			ParkAllowedMinPosition:     2,
			ParkResumeRetryIntervalMin: 1,
			ParkResumeMaxAttempts:      3,
			ParkResumeMenuTimeoutSec:   15,
			ParkResumeDialTimeoutSec:   20,
		}
		if err := DB.Create(&defaultQueue).Error; err == nil {
			queueMembers := []QueueMember{
				{QueueID: defaultQueue.ID, Member: "Lorraine Robertson", Extension: "1002", Penalty: 0, Disabled: false},
				{QueueID: defaultQueue.ID, Member: "Sharon Collinson", Extension: "1003", Penalty: 0, Disabled: false},
				{QueueID: defaultQueue.ID, Member: "Julie Schubeler", Extension: "1004", Penalty: 0, Disabled: false},
				{QueueID: defaultQueue.ID, Member: "Amy Wood", Extension: "1005", Penalty: 0, Disabled: false},
				{QueueID: defaultQueue.ID, Member: "Channelle Brown", Extension: "1006", Penalty: 0, Disabled: false},
			}
			for _, qm := range queueMembers {
				_ = DB.Create(&qm)
			}
		}
	}

	// Ensure realistic sample CDR records exist for drill-down reporting & tracking
	var cdrCount int64
	DB.Model(&CDR{}).Count(&cdrCount)
	if cdrCount < 8 {
		var firstTenant Tenant
		if err := DB.First(&firstTenant).Error; err == nil {
			now := time.Now()
			sampleCDRs := []CDR{
				{
					Base:           Base{CreatedAt: now.Add(-25 * time.Minute)},
					TenantID:       firstTenant.ID,
					CallerIDName:   "Enterprise Partner",
					CallerIDNumber: "+15550199456",
					Destination:    "1001",
					Context:        "public",
					Duration:       245,
					Billsec:        230,
					HangupCause:    "NORMAL_CLEARING",
					UUID:           "a1b2c3d4-e5f6-7890-abcd-111111111111",
					Transcription:  "Inbound inquiry regarding telecom trunk capacity and SIP SLA parameters.",
				},
				{
					Base:           Base{CreatedAt: now.Add(-1 * time.Hour - 15*time.Minute)},
					TenantID:       firstTenant.ID,
					CallerIDName:   "Alice Smith",
					CallerIDNumber: "1001",
					Destination:    "+18005550199",
					Context:        "outbound",
					Duration:       180,
					Billsec:        165,
					HangupCause:    "NORMAL_CLEARING",
					UUID:           "a1b2c3d4-e5f6-7890-abcd-222222222222",
				},
				{
					Base:           Base{CreatedAt: now.Add(-2 * time.Hour)},
					TenantID:       firstTenant.ID,
					CallerIDName:   "Customer VIP",
					CallerIDNumber: "+442079460123",
					Destination:    "1002",
					Context:        "public",
					Duration:       420,
					Billsec:        405,
					HangupCause:    "NORMAL_CLEARING",
					UUID:           "a1b2c3d4-e5f6-7890-abcd-333333333333",
					Transcription:  "Client confirmed invoice renewal and asked for priority routing schedule.",
				},
				{
					Base:           Base{CreatedAt: now.Add(-3 * time.Hour - 40*time.Minute)},
					TenantID:       firstTenant.ID,
					CallerIDName:   "Unknown Caller",
					CallerIDNumber: "+15550144321",
					Destination:    "1001",
					Context:        "public",
					Duration:       25,
					Billsec:        0,
					HangupCause:    "NO_ANSWER",
					UUID:           "a1b2c3d4-e5f6-7890-abcd-444444444444",
				},
				{
					Base:           Base{CreatedAt: now.Add(-4 * time.Hour)},
					TenantID:       firstTenant.ID,
					CallerIDName:   "Bob Jones",
					CallerIDNumber: "1002",
					Destination:    "1001",
					Context:        "internal",
					Duration:       95,
					Billsec:        90,
					HangupCause:    "NORMAL_CLEARING",
					UUID:           "a1b2c3d4-e5f6-7890-abcd-555555555555",
				},
				{
					Base:           Base{CreatedAt: now.Add(-24 * time.Hour - 2*time.Hour)},
					TenantID:       firstTenant.ID,
					CallerIDName:   "Carrier NOC",
					CallerIDNumber: "+18005550000",
					Destination:    "600",
					Context:        "public",
					Duration:       310,
					Billsec:        295,
					HangupCause:    "NORMAL_CLEARING",
					UUID:           "a1b2c3d4-e5f6-7890-abcd-666666666666",
				},
				{
					Base:           Base{CreatedAt: now.Add(-26 * time.Hour)},
					TenantID:       firstTenant.ID,
					CallerIDName:   "Telephony User",
					CallerIDNumber: "1003",
					Destination:    "*97",
					Context:        "internal",
					Duration:       45,
					Billsec:        40,
					HangupCause:    "NORMAL_CLEARING",
					UUID:           "a1b2c3d4-e5f6-7890-abcd-777777777777",
				},
				{
					Base:           Base{CreatedAt: now.Add(-48 * time.Hour)},
					TenantID:       firstTenant.ID,
					CallerIDName:   "Sales Lead",
					CallerIDNumber: "+15550122999",
					Destination:    "1001",
					Context:        "public",
					Duration:       15,
					Billsec:        0,
					HangupCause:    "USER_BUSY",
					UUID:           "a1b2c3d4-e5f6-7890-abcd-888888888888",
				},
			}
			for _, c := range sampleCDRs {
				_ = DB.Create(&c)
			}
		}
	}
}
