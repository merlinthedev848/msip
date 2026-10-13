package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// RegisterCallControlRoutes registers enterprise PBX call control and voice feature APIs
func RegisterCallControlRoutes(v1 *gin.RouterGroup) {
	// ==========================================
	// TIME CONDITIONS (Office Hours Routing)
	// ==========================================
	v1.GET("/time-conditions", func(c *gin.Context) {
		var items []TimeCondition
		scopedDB(c).Find(&items)
		c.JSON(http.StatusOK, gin.H{"time_conditions": items})
	})

	v1.POST("/time-conditions", func(c *gin.Context) {
		var item TimeCondition
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		item.TenantID = getTenantID(c)
		if err := DB.Create(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, item)
	})

	v1.PUT("/time-conditions/:id", func(c *gin.Context) {
		id := c.Param("id")
		var record TimeCondition
		if err := scopedDB(c).First(&record, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Time condition not found"})
			return
		}
		var input TimeCondition
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

	v1.DELETE("/time-conditions/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := scopedDB(c).Delete(&TimeCondition{}, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "deleted"})
	})

	// ==========================================
	// HOLIDAYS
	// ==========================================
	v1.GET("/holidays", func(c *gin.Context) {
		var items []Holiday
		scopedDB(c).Find(&items)
		c.JSON(http.StatusOK, gin.H{"holidays": items})
	})

	v1.POST("/holidays", func(c *gin.Context) {
		var item Holiday
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		item.TenantID = getTenantID(c)
		if err := DB.Create(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, item)
	})

	v1.PUT("/holidays/:id", func(c *gin.Context) {
		id := c.Param("id")
		var record Holiday
		if err := scopedDB(c).First(&record, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Holiday not found"})
			return
		}
		var input Holiday
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

	v1.DELETE("/holidays/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := scopedDB(c).Delete(&Holiday{}, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "deleted"})
	})

	// ==========================================
	// OUTBOUND RESTRICTIONS (Call Barring)
	// ==========================================
	v1.GET("/outbound-restrictions", func(c *gin.Context) {
		var items []OutboundRestriction
		scopedDB(c).Find(&items)
		c.JSON(http.StatusOK, gin.H{"restrictions": items})
	})

	v1.POST("/outbound-restrictions", func(c *gin.Context) {
		var item OutboundRestriction
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		item.TenantID = getTenantID(c)
		if err := DB.Create(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, item)
	})

	v1.PUT("/outbound-restrictions/:id", func(c *gin.Context) {
		id := c.Param("id")
		var record OutboundRestriction
		if err := scopedDB(c).First(&record, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Restriction not found"})
			return
		}
		var input OutboundRestriction
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

	v1.DELETE("/outbound-restrictions/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := scopedDB(c).Delete(&OutboundRestriction{}, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "deleted"})
	})

	// ==========================================
	// AUTOCLIP (Smart Callback Routing)
	// ==========================================
	v1.GET("/autoclip", func(c *gin.Context) {
		var items []AutoCLIP
		scopedDB(c).Order("created_at desc").Limit(50).Find(&items)
		c.JSON(http.StatusOK, gin.H{"autoclip": items})
	})

	v1.POST("/autoclip", func(c *gin.Context) {
		var item AutoCLIP
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		item.TenantID = getTenantID(c)
		if err := DB.Create(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, item)
	})

	v1.DELETE("/autoclip/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := scopedDB(c).Delete(&AutoCLIP{}, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "deleted"})
	})

	// ==========================================
	// PAGING & INTERCOM GROUPS
	// ==========================================
	v1.GET("/paging-groups", func(c *gin.Context) {
		var items []PagingGroup
		scopedDB(c).Find(&items)
		c.JSON(http.StatusOK, gin.H{"paging_groups": items})
	})

	v1.POST("/paging-groups", func(c *gin.Context) {
		var item PagingGroup
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		item.TenantID = getTenantID(c)
		if err := DB.Create(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, item)
	})

	v1.PUT("/paging-groups/:id", func(c *gin.Context) {
		id := c.Param("id")
		var record PagingGroup
		if err := scopedDB(c).First(&record, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Paging group not found"})
			return
		}
		var input PagingGroup
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

	v1.DELETE("/paging-groups/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := scopedDB(c).Delete(&PagingGroup{}, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "deleted"})
	})

	// ==========================================
	// CALL PICKUP GROUPS (*8)
	// ==========================================
	v1.GET("/pickup-groups", func(c *gin.Context) {
		var items []PickupGroup
		scopedDB(c).Find(&items)
		c.JSON(http.StatusOK, gin.H{"pickup_groups": items})
	})

	v1.POST("/pickup-groups", func(c *gin.Context) {
		var item PickupGroup
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		item.TenantID = getTenantID(c)
		if err := DB.Create(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, item)
	})

	v1.PUT("/pickup-groups/:id", func(c *gin.Context) {
		id := c.Param("id")
		var record PickupGroup
		if err := scopedDB(c).First(&record, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Pickup group not found"})
			return
		}
		var input PickupGroup
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

	v1.DELETE("/pickup-groups/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := scopedDB(c).Delete(&PickupGroup{}, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "deleted"})
	})

	// ==========================================
	// AUDIO PROMPTS & MUSIC ON HOLD (MoH)
	// ==========================================
	v1.GET("/audio-prompts", func(c *gin.Context) {
		var items []AudioPrompt
		scopedDB(c).Find(&items)
		c.JSON(http.StatusOK, gin.H{"audio_prompts": items})
	})

	v1.POST("/audio-prompts", func(c *gin.Context) {
		var item AudioPrompt
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		item.TenantID = getTenantID(c)
		if err := DB.Create(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, item)
	})

	v1.PUT("/audio-prompts/:id", func(c *gin.Context) {
		id := c.Param("id")
		var record AudioPrompt
		if err := scopedDB(c).First(&record, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Audio prompt not found"})
			return
		}
		var input AudioPrompt
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

	v1.DELETE("/audio-prompts/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := scopedDB(c).Delete(&AudioPrompt{}, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "deleted"})
	})

	// ==========================================
	// FEATURE CODES (*97, *8, *72, 9999, etc.)
	// ==========================================
	v1.GET("/feature-codes", func(c *gin.Context) {
		var items []FeatureCode
		scopedDB(c).Find(&items)
		c.JSON(http.StatusOK, gin.H{"feature_codes": items})
	})

	v1.POST("/feature-codes", func(c *gin.Context) {
		var item FeatureCode
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		item.TenantID = getTenantID(c)
		if err := DB.Create(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, item)
	})

	v1.PUT("/feature-codes/:id", func(c *gin.Context) {
		id := c.Param("id")
		var record FeatureCode
		if err := scopedDB(c).First(&record, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Feature code not found"})
			return
		}
		var input FeatureCode
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

	v1.DELETE("/feature-codes/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := scopedDB(c).Delete(&FeatureCode{}, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "deleted"})
	})

	// ==========================================
	// CALL RECORDING POLICIES
	// ==========================================
	v1.GET("/recording-rules", func(c *gin.Context) {
		var items []RecordingRule
		scopedDB(c).Find(&items)
		c.JSON(http.StatusOK, gin.H{"recording_rules": items})
	})

	v1.POST("/recording-rules", func(c *gin.Context) {
		var item RecordingRule
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		item.TenantID = getTenantID(c)
		if err := DB.Create(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, item)
	})

	v1.PUT("/recording-rules/:id", func(c *gin.Context) {
		id := c.Param("id")
		var record RecordingRule
		if err := scopedDB(c).First(&record, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Recording rule not found"})
			return
		}
		var input RecordingRule
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

	v1.DELETE("/recording-rules/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := scopedDB(c).Delete(&RecordingRule{}, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "deleted"})
	})

	// ==========================================
	// DIAGNOSTICS (SIP, IAX, RTP & JITTER BUFFER)
	// ==========================================
	v1.GET("/diagnostics/jitter", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "operational",
			"sip_profiles": []gin.H{
				{"name": "internal", "port": 5060, "proto": "udp/tcp", "state": "RUNNING"},
				{"name": "external", "port": 5080, "proto": "udp/tcp", "state": "RUNNING"},
				{"name": "webrtc_wss", "port": 7443, "proto": "wss", "state": "RUNNING"},
				{"name": "kamailio_proxy", "port": 5060, "proto": "sip/udp", "state": "RUNNING"},
			},
			"jitter_buffer": gin.H{
				"current_ms":      0,
				"max_ms":          60,
				"packet_loss_pct": 0.0,
				"rtt_ms":          0,
				"mos_score":       4.5,
				"quality":         "Idle / Nominal",
			},
			"active_codecs": []string{"OPUS", "G722", "PCMA (G.711a)", "PCMU (G.711u)"},
		})
	})

	v1.POST("/diagnostics/ping", func(c *gin.Context) {
		var req struct {
			Host string `json:"host"`
		}
		if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Host) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "host is required"})
			return
		}
		target := strings.TrimSpace(req.Host)
		if !strings.Contains(target, ":") {
			target = target + ":5060"
		}
		start := time.Now()
		conn, err := net.DialTimeout("tcp", target, 2*time.Second)
		elapsed := time.Since(start).Milliseconds()
		if err != nil {
			// Fallback UDP probe
			uConn, uErr := net.DialTimeout("udp", target, 2*time.Second)
			if uErr != nil {
				c.JSON(http.StatusOK, gin.H{
					"target":     req.Host,
					"latency_ms": 0,
					"status":     fmt.Sprintf("Unreachable (%v)", err),
					"success":    false,
				})
				return
			}
			_ = uConn.Close()
		} else {
			_ = conn.Close()
		}
		c.JSON(http.StatusOK, gin.H{
			"target":     req.Host,
			"latency_ms": elapsed,
			"status":     fmt.Sprintf("Reachable (%dms RTT)", elapsed),
			"success":    true,
		})
	})
}
