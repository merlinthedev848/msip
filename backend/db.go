package main

import (
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Fallback for local docker-compose
		dsn = "host=localhost user=msip_user password=msip_secure dbname=msip_pbx port=5432 sslmode=disable TimeZone=UTC"
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Ensure the pgcrypto extension is created for gen_random_uuid()
	DB.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto;")

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
	)
	if err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}
	
	// Seed a default tenant if none exists
	var count int64
	DB.Model(&Tenant{}).Count(&count)
	if count == 0 {
		defaultTenant := Tenant{Name: "Default Organization"}
		DB.Create(&defaultTenant)
		log.Printf("Seeded default tenant: %s", defaultTenant.ID.String())

		// Seed a default admin user
		hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		defaultUser := User{
			TenantID:     defaultTenant.ID,
			Email:        "admin@pbx.local",
			PasswordHash: string(hash),
			Role:         "superadmin",
		}
		DB.Create(&defaultUser)
		log.Printf("Seeded default admin user: admin@pbx.local")

		// Seed a default test extension: 1001 / testpass
		defaultExt := Extension{
			TenantID:        defaultTenant.ID,
			ExtensionNumber: "1001",
			PasswordHash:    "testpass",
			IsActive:        true,
		}
		DB.Create(&defaultExt)
		log.Printf("Seeded default extension: 1001")
	}
}
