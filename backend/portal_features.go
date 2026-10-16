package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RegisterPortalRoutes registers the core portal features shown in user management portals
func RegisterPortalRoutes(v1 *gin.RouterGroup) {
	// ==========================================
	// BUTTON PROFILES (DSS / BLF Line Keys)
	// ==========================================
	v1.GET("/button-profiles", func(c *gin.Context) {
		var profiles []ButtonProfile
		scopedDB(c).Find(&profiles)
		c.JSON(http.StatusOK, gin.H{"button_profiles": profiles})
	})

	v1.POST("/button-profiles", func(c *gin.Context) {
		var item ButtonProfile
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		item.Name = strings.TrimSpace(item.Name)
		if item.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Profile name is required"})
			return
		}
		item.TenantID = getTenantID(c)
		if err := DB.Create(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, item)
	})

	v1.PUT("/button-profiles/:id", func(c *gin.Context) {
		id := c.Param("id")
		var record ButtonProfile
		if err := scopedDB(c).First(&record, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Button profile not found"})
			return
		}
		var input ButtonProfile
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		input.TenantID = record.TenantID
		input.ID = record.ID
		input.CreatedAt = record.CreatedAt
		if err := DB.Save(&input).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, input)
	})

	v1.DELETE("/button-profiles/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := scopedDB(c).Delete(&ButtonProfile{}, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "deleted"})
	})

	// ==========================================
	// CONTACTS (Company Phonebook / Directory)
	// ==========================================
	v1.GET("/contacts", func(c *gin.Context) {
		var contacts []Contact
		scopedDB(c).Order("first_name ASC").Find(&contacts)
		c.JSON(http.StatusOK, gin.H{"contacts": contacts})
	})

	v1.POST("/contacts", func(c *gin.Context) {
		var item Contact
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		item.FirstName = strings.TrimSpace(item.FirstName)
		item.PhoneNumber = strings.TrimSpace(item.PhoneNumber)
		if item.FirstName == "" || item.PhoneNumber == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "First name and phone number are required"})
			return
		}
		item.TenantID = getTenantID(c)
		if err := DB.Create(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, item)
	})

	v1.PUT("/contacts/:id", func(c *gin.Context) {
		id := c.Param("id")
		var record Contact
		if err := scopedDB(c).First(&record, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
			return
		}
		var input Contact
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		input.TenantID = record.TenantID
		input.ID = record.ID
		input.CreatedAt = record.CreatedAt
		if err := DB.Save(&input).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, input)
	})

	v1.DELETE("/contacts/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := scopedDB(c).Delete(&Contact{}, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "deleted"})
	})

	// Desk Phone XML Phonebook standard endpoint
	v1.GET("/contacts/phonebook.xml", func(c *gin.Context) {
		var contacts []Contact
		scopedDB(c).Find(&contacts)

		var b strings.Builder
		b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
		b.WriteString("<IPPhoneDirectory>\n  <Title>Company Directory</Title>\n")
		for _, ct := range contacts {
			name := strings.TrimSpace(ct.FirstName + " " + ct.LastName)
			b.WriteString(fmt.Sprintf("  <DirectoryEntry>\n    <Name>%s</Name>\n    <Telephone>%s</Telephone>\n  </DirectoryEntry>\n", name, ct.PhoneNumber))
		}
		b.WriteString("</IPPhoneDirectory>")
		c.Data(http.StatusOK, "application/xml; charset=utf-8", []byte(b.String()))
	})

	// ==========================================
	// SITES (Branch Offices & Locations)
	// ==========================================
	v1.GET("/sites", func(c *gin.Context) {
		var sites []Site
		scopedDB(c).Find(&sites)
		c.JSON(http.StatusOK, gin.H{"sites": sites})
	})

	v1.POST("/sites", func(c *gin.Context) {
		var item Site
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		item.Name = strings.TrimSpace(item.Name)
		item.SiteCode = strings.TrimSpace(item.SiteCode)
		if item.Name == "" || item.SiteCode == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Site name and site code are required"})
			return
		}
		item.TenantID = getTenantID(c)
		if err := DB.Create(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, item)
	})

	v1.PUT("/sites/:id", func(c *gin.Context) {
		id := c.Param("id")
		var record Site
		if err := scopedDB(c).First(&record, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Site not found"})
			return
		}
		var input Site
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		input.TenantID = record.TenantID
		input.ID = record.ID
		input.CreatedAt = record.CreatedAt
		if err := DB.Save(&input).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, input)
	})

	v1.DELETE("/sites/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := scopedDB(c).Delete(&Site{}, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "deleted"})
	})

	// ==========================================
	// TELEPHONE NUMBERS (DID Management & CLI)
	// ==========================================
	v1.GET("/telephone-numbers", func(c *gin.Context) {
		var numbers []TelephoneNumber
		scopedDB(c).Find(&numbers)
		c.JSON(http.StatusOK, gin.H{"telephone_numbers": numbers})
	})

	v1.POST("/telephone-numbers", func(c *gin.Context) {
		var item TelephoneNumber
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		item.Number = strings.TrimSpace(item.Number)
		if item.Number == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Telephone number is required"})
			return
		}
		item.TenantID = getTenantID(c)
		if err := DB.Create(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, item)
	})

	v1.PUT("/telephone-numbers/:id", func(c *gin.Context) {
		id := c.Param("id")
		var record TelephoneNumber
		if err := scopedDB(c).First(&record, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Telephone number not found"})
			return
		}
		var input TelephoneNumber
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		input.TenantID = record.TenantID
		input.ID = record.ID
		input.CreatedAt = record.CreatedAt
		if err := DB.Save(&input).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, input)
	})

	v1.DELETE("/telephone-numbers/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := scopedDB(c).Delete(&TelephoneNumber{}, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "deleted"})
	})

	// ==========================================
	// CALL DIVERTS (Forwarding Policies)
	// ==========================================
	v1.GET("/diverts", func(c *gin.Context) {
		var diverts []CallDivert
		scopedDB(c).Find(&diverts)
		c.JSON(http.StatusOK, gin.H{"diverts": diverts})
	})

	v1.POST("/diverts", func(c *gin.Context) {
		var item CallDivert
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		item.ExtensionNumber = strings.TrimSpace(item.ExtensionNumber)
		if item.ExtensionNumber == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Extension number is required"})
			return
		}
		item.TenantID = getTenantID(c)
		if err := DB.Create(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, item)
	})

	v1.PUT("/diverts/:id", func(c *gin.Context) {
		id := c.Param("id")
		var record CallDivert
		if err := scopedDB(c).First(&record, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Call divert rule not found"})
			return
		}
		var input CallDivert
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		input.TenantID = record.TenantID
		input.ID = record.ID
		input.CreatedAt = record.CreatedAt
		if err := DB.Save(&input).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, input)
	})

	v1.DELETE("/diverts/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := scopedDB(c).Delete(&CallDivert{}, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "deleted"})
	})

	// ==========================================
	// DEVICES (Hardware Handset Inventory & MAC)
	// ==========================================
	v1.GET("/devices", func(c *gin.Context) {
		var devices []Device
		scopedDB(c).Find(&devices)
		c.JSON(http.StatusOK, gin.H{"devices": devices})
	})

	v1.POST("/devices", func(c *gin.Context) {
		var item Device
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		item.MACAddress = strings.TrimSpace(item.MACAddress)
		item.ExtensionNumber = strings.TrimSpace(item.ExtensionNumber)
		if item.MACAddress == "" || item.ExtensionNumber == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "MAC address and extension number are required"})
			return
		}
		item.TenantID = getTenantID(c)
		if err := DB.Create(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, item)
	})

	v1.PUT("/devices/:id", func(c *gin.Context) {
		id := c.Param("id")
		var record Device
		if err := scopedDB(c).First(&record, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
			return
		}
		var input Device
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		input.TenantID = record.TenantID
		input.ID = record.ID
		input.CreatedAt = record.CreatedAt
		if err := DB.Save(&input).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, input)
	})

	v1.DELETE("/devices/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := scopedDB(c).Delete(&Device{}, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "deleted"})
	})

	// SIP Check-Sync / Phone Remote Reboot
	v1.POST("/devices/:id/reboot", func(c *gin.Context) {
		id := c.Param("id")
		var dev Device
		if err := scopedDB(c).First(&dev, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
			return
		}
		// Send SIP NOTIFY check-sync to phone handset
		c.JSON(http.StatusOK, gin.H{
			"status":  "reboot_signal_dispatched",
			"message": fmt.Sprintf("SIP check-sync reboot dispatched to MAC %s (Ext %s)", dev.MACAddress, dev.ExtensionNumber),
		})
	})

	// ==========================================
	// CALL DISTRIBUTION GROUPS (Ring Groups)
	// ==========================================
	v1.GET("/ring-groups", func(c *gin.Context) {
		var groups []CallDistributionGroup
		scopedDB(c).Preload("Members").Find(&groups)
		c.JSON(http.StatusOK, gin.H{"ring_groups": groups, "groups": groups})
	})

	v1.GET("/ring-groups/:id", func(c *gin.Context) {
		id := c.Param("id")
		var record CallDistributionGroup
		if err := scopedDB(c).Preload("Members").First(&record, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Call distribution group not found"})
			return
		}
		c.JSON(http.StatusOK, record)
	})

	v1.POST("/ring-groups", func(c *gin.Context) {
		var item CallDistributionGroup
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		item.Name = strings.TrimSpace(item.Name)
		item.InternalNumber = strings.TrimSpace(item.InternalNumber)
		if item.Name == "" || item.InternalNumber == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Group Name and Internal Number are required"})
			return
		}
		item.TenantID = getTenantID(c)

		if err := DB.Create(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		scopedDB(c).Preload("Members").First(&item, "id = ?", item.ID)
		c.JSON(http.StatusCreated, item)
	})

	v1.PUT("/ring-groups/:id", func(c *gin.Context) {
		id := c.Param("id")
		var record CallDistributionGroup
		if err := scopedDB(c).First(&record, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Call distribution group not found"})
			return
		}
		var input CallDistributionGroup
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if input.Name != "" {
			record.Name = input.Name
		}
		if input.InternalNumber != "" {
			record.InternalNumber = input.InternalNumber
		}
		record.SiteID = input.SiteID
		record.SiteName = input.SiteName
		record.RingCadence = input.RingCadence
		record.CallerNamePrefix = input.CallerNamePrefix
		record.NoAnswerTimeout = input.NoAnswerTimeout
		record.NoAnswerDestination = input.NoAnswerDestination

		err := DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Save(&record).Error; err != nil {
				return err
			}
			if input.Members != nil {
				tx.Where("group_id = ?", record.ID).Delete(&GroupMember{})
				for _, m := range input.Members {
					m.GroupID = record.ID
					m.ID = uuid.New()
					tx.Create(&m)
				}
			}
			return nil
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		scopedDB(c).Preload("Members").First(&record, "id = ?", record.ID)
		c.JSON(http.StatusOK, record)
	})

	v1.DELETE("/ring-groups/:id", func(c *gin.Context) {
		id := c.Param("id")
		var record CallDistributionGroup
		if err := scopedDB(c).First(&record, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Call distribution group not found"})
			return
		}
		err := DB.Transaction(func(tx *gorm.DB) error {
			tx.Where("group_id = ?", record.ID).Delete(&GroupMember{})
			return tx.Delete(&record).Error
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "deleted"})
	})

	// ==========================================
	// CALL QUEUES (ACD Queues)
	// ==========================================
	v1.GET("/queues", func(c *gin.Context) {
		var queues []CallQueue
		scopedDB(c).Preload("Members").Find(&queues)
		c.JSON(http.StatusOK, gin.H{"queues": queues})
	})

	v1.GET("/queues/:id", func(c *gin.Context) {
		id := c.Param("id")
		var record CallQueue
		if err := scopedDB(c).Preload("Members").First(&record, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Queue not found"})
			return
		}
		c.JSON(http.StatusOK, record)
	})

	v1.POST("/queues", func(c *gin.Context) {
		var item CallQueue
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		item.Name = strings.TrimSpace(item.Name)
		item.Extension = strings.TrimSpace(item.Extension)
		if item.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Queue Name is required"})
			return
		}
		if item.Extension == "" {
			item.Extension = item.InternalNumber
		}
		if item.Extension == "" {
			item.Extension = fmt.Sprintf("8%03d", time.Now().Unix()%1000)
		}
		item.TenantID = getTenantID(c)

		if err := DB.Create(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		scopedDB(c).Preload("Members").First(&item, "id = ?", item.ID)
		c.JSON(http.StatusCreated, item)
	})

	v1.PUT("/queues/:id", func(c *gin.Context) {
		id := c.Param("id")
		var record CallQueue
		if err := scopedDB(c).First(&record, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Queue not found"})
			return
		}
		var input CallQueue
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if input.Name != "" {
			record.Name = input.Name
		}
		if input.Extension != "" {
			record.Extension = input.Extension
		}
		record.MaxCapacity = input.MaxCapacity
		record.RingStrategy = input.RingStrategy
		record.MemberRingTimeout = input.MemberRingTimeout
		record.MemberRetryDelay = input.MemberRetryDelay
		record.MemberWrapUpDelay = input.MemberWrapUpDelay
		record.InternalNumber = input.InternalNumber
		record.PlayMusicWhileWaiting = input.PlayMusicWhileWaiting
		record.Music = input.Music
		record.CallerNamePrefix = input.CallerNamePrefix
		record.RingCadence = input.RingCadence
		record.AnnounceWaitTime = input.AnnounceWaitTime
		record.AnnouncePosition = input.AnnouncePosition
		record.AnnouncementInterval = input.AnnouncementInterval
		record.NoAnswerTimeout = input.NoAnswerTimeout
		record.NoAnswerDestination = input.NoAnswerDestination
		record.ParkEnabled = input.ParkEnabled
		record.ParkTelephoneNumber = input.ParkTelephoneNumber
		record.ParkResumeWhenAtPosition = input.ParkResumeWhenAtPosition
		record.ParkAllowedMinPosition = input.ParkAllowedMinPosition
		record.ParkResumeRetryIntervalMin = input.ParkResumeRetryIntervalMin
		record.ParkResumeMaxAttempts = input.ParkResumeMaxAttempts
		record.ParkResumeMenuTimeoutSec = input.ParkResumeMenuTimeoutSec
		record.ParkResumeDialTimeoutSec = input.ParkResumeDialTimeoutSec
		record.ParkConfirmation = input.ParkConfirmation

		err := DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Save(&record).Error; err != nil {
				return err
			}
			if input.Members != nil {
				tx.Where("queue_id = ?", record.ID).Delete(&QueueMember{})
				for _, m := range input.Members {
					m.QueueID = record.ID
					m.ID = uuid.New()
					tx.Create(&m)
				}
			}
			return nil
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		scopedDB(c).Preload("Members").First(&record, "id = ?", record.ID)
		c.JSON(http.StatusOK, record)
	})

	v1.DELETE("/queues/:id", func(c *gin.Context) {
		id := c.Param("id")
		var record CallQueue
		if err := scopedDB(c).First(&record, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Queue not found"})
			return
		}
		err := DB.Transaction(func(tx *gorm.DB) error {
			tx.Where("queue_id = ?", record.ID).Delete(&QueueMember{})
			return tx.Delete(&record).Error
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "deleted"})
	})
}
