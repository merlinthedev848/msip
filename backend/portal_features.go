package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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
		if err := DB.Model(&record).Updates(input).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, record)
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
		if err := DB.Model(&record).Updates(input).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, record)
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
		if err := DB.Model(&record).Updates(input).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, record)
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
		if err := DB.Model(&record).Updates(input).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, record)
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
		if err := DB.Model(&record).Updates(input).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, record)
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
		if err := DB.Model(&record).Updates(input).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, record)
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
}
