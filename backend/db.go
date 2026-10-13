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
	}
}
