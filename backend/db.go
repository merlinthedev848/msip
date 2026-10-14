package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
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
	maxAttempts := 15
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		DB, err = gorm.Open(dialector, &gorm.Config{
			Logger: logger.Default.LogMode(logger.Warn),
		})
		if err == nil {
			break
		}
		if attempt < maxAttempts {
			log.Printf("Waiting for database connection (attempt %d/%d): %v", attempt, maxAttempts, err)
			time.Sleep(2 * time.Second)
		}
	}
	if err != nil {
		log.Fatalf("Failed to connect to database after %d attempts: %v", maxAttempts, err)
	}

	sqlDB, err := DB.DB()
	if err == nil {
		sqlDB.SetMaxIdleConns(25)
		sqlDB.SetMaxOpenConns(150)
		sqlDB.SetConnMaxLifetime(time.Hour)
		sqlDB.SetConnMaxIdleTime(15 * time.Minute)
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
		&SIPSettings{},
		&AIAssistant{},
		&AIBooking{},
		&AICallSession{},
		&AITranscriptionSettings{},
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

		// Seed initial SIP & RTP settings
		defaultSIP := DefaultSIPSettings(defaultTenant.ID)
		if err := DB.Create(&defaultSIP).Error; err != nil {
			log.Printf("Error seeding default SIP settings: %v", err)
		} else {
			log.Printf("Seeded default SIP & RTP settings for tenant: %s", defaultTenant.ID.String())
		}

		// Seed initial AI & Transcription Settings
		defaultAISettings := DefaultAITranscriptionSettings(defaultTenant.ID)
		_ = DB.Create(&defaultAISettings)

		// Seed default AI Receptionist & Appointment Booking Assistant
		defaultAssistant := AIAssistant{
			Base:               Base{ID: uuid.New()},
			TenantRef:          TenantRef{TenantID: defaultTenant.ID},
			Name:               "Sarah - AI Virtual Receptionist",
			Extension:          "890",
			RoleType:           "receptionist",
			Greeting:           "Hello! Thank you for calling MSIP Telephony. I am Sarah, your AI virtual receptionist. I can help book appointments, answer questions, or transfer your call to an extension. How may I assist you today?",
			SystemPrompt:       "You are Sarah, a warm, professional virtual receptionist for an enterprise company using MSIP PBX. You speak clearly and concisely. You can answer business questions (hours: Mon-Fri 8am-6pm), schedule customer appointments, or transfer calls to extensions (1001 for Technical Support, 1002 for Billing). If the caller wants to book an appointment, collect their name, desired date, time, and service type, and confirm the details.",
			LLMProvider:        "openai",
			LLMModel:           "gpt-4o",
			VoiceProvider:      "openai",
			VoiceName:          "alloy",
			Language:           "en-US",
			EnableBooking:      true,
			EnableTransfer:     true,
			TransferTarget:     "1001",
			EnableVoicemail:    true,
			MaxCallDurationSec: 600,
			IsActive:           true,
		}
		if err := DB.Create(&defaultAssistant).Error; err != nil {
			log.Printf("Error seeding default AI assistant: %v", err)
		} else {
			log.Printf("Seeded default AI assistant: %s (Ext %s)", defaultAssistant.Name, defaultAssistant.Extension)
		}
	}
}
