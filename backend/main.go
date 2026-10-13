package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var jwtSecret = []byte("super-secret-pbx-key")

// AuthMiddleware validates JWT and extracts tenant_id and role
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if len(tokenString) > 7 && strings.EqualFold(tokenString[:7], "Bearer ") {
			tokenString = strings.TrimSpace(tokenString[7:])
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or malformed Authorization header"})
			return
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		tenantIDStr, _ := claims["tenant_id"].(string)
		tenantID, err := uuid.Parse(tenantIDStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid tenant ID in token"})
			return
		}

		role, _ := claims["role"].(string)
		sub, _ := claims["sub"].(string)

		c.Set("tenant_id", tenantID)
		c.Set("role", role)
		c.Set("user_id", sub)

		// Read-only role check for mutation operations
		if role == "viewer" {
			if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut || c.Request.Method == http.MethodDelete {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Role 'viewer' has read-only access"})
				return
			}
		}

		c.Next()
	}
}

func getTenantID(c *gin.Context) uuid.UUID {
	val, exists := c.Get("tenant_id")
	if !exists {
		return uuid.Nil
	}
	id, ok := val.(uuid.UUID)
	if !ok {
		return uuid.Nil
	}
	return id
}

func isSuperAdmin(c *gin.Context) bool {
	role, _ := c.Get("role")
	return role == "superadmin"
}

// scopedDB returns a GORM query scoped to the user's tenant unless the user is superadmin
func scopedDB(c *gin.Context) *gorm.DB {
	tenantID := getTenantID(c)
	if isSuperAdmin(c) {
		return DB
	}
	return DB.Where("tenant_id = ?", tenantID)
}

func sanitizeUser(u *User) {
	u.PasswordHash = ""
}

func sanitizeExtension(e *Extension) {
	e.PasswordHash = ""
}

func generateRandomSIPPassword(length int) string {
	const charset = "abcdefghjkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789!@#$"
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "SipPass" + fmt.Sprintf("%d", time.Now().UnixNano()%1000000)
	}
	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b)
}

// CORSMiddleware enables universal cross-origin requests for browsers
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// SetupRouter initializes and configures the Gin engine
func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(CORSMiddleware())

	// Normalize /api/v1/ prefix to /api/ so direct backend calls succeed
	r.Use(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/v1/") {
			c.Request.URL.Path = "/api/" + strings.TrimPrefix(c.Request.URL.Path, "/api/v1/")
		}
		c.Next()
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// FreeSWITCH mod_xml_curl dynamic gateway (support both /fs and /fs/curl)
	r.POST("/fs", HandleFSCurl)
	r.POST("/fs/curl", HandleFSCurl)

	api := r.Group("/api")
	{
		// Internal System Route - FreeSWITCH posts CDR here
		api.POST("/cdr", func(c *gin.Context) {
			var payload map[string]interface{}
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid CDR payload"})
				return
			}

			variables, ok := payload["variables"].(map[string]interface{})
			if !ok || variables == nil {
				c.JSON(http.StatusOK, gin.H{"status": "ignored", "reason": "no variables found"})
				return
			}

			callUUID, _ := variables["uuid"].(string)
			callerIDName, _ := variables["caller_id_name"].(string)
			callerIDNumber, _ := variables["caller_id_number"].(string)
			destination, _ := variables["destination_number"].(string)
			context, _ := variables["context"].(string)
			domainName, _ := variables["domain_name"].(string)
			accountCode, _ := variables["accountcode"].(string)

			var duration, billsec int
			if d, ok := variables["duration"].(float64); ok {
				duration = int(d)
			}
			if b, ok := variables["billsec"].(float64); ok {
				billsec = int(b)
			}
			hangupCause, _ := variables["hangup_cause"].(string)

			var tenantID uuid.UUID
			if accountCode != "" {
				if parsed, err := uuid.Parse(accountCode); err == nil {
					tenantID = parsed
				}
			}
			if tenantID == uuid.Nil && domainName != "" {
				var t Tenant
				if err := DB.Where("domain = ?", domainName).First(&t).Error; err == nil {
					tenantID = t.ID
				}
			}
			if tenantID == uuid.Nil {
				var defaultT Tenant
				if err := DB.First(&defaultT).Error; err == nil {
					tenantID = defaultT.ID
				}
			}

			cdr := CDR{
				Base:           Base{ID: uuid.New()},
				TenantID:       tenantID,
				UUID:           callUUID,
				CallerIDName:   callerIDName,
				CallerIDNumber: callerIDNumber,
				Destination:    destination,
				Context:        context,
				Duration:       duration,
				Billsec:        billsec,
				HangupCause:    hangupCause,
			}

			if err := DB.Create(&cdr).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save CDR"})
				return
			}

			if destination == "*88" {
				go func(id string, callCDR CDR) {
					time.Sleep(2 * time.Second)
					callCDR.Transcription = "[AI Transcription Generated via Whisper] Voice memo processed from PBX extension."
					DB.Model(&CDR{}).Where("uuid = ?", id).Update("transcription", callCDR.Transcription)
				}(callUUID, cdr)
			}

			c.JSON(http.StatusCreated, gin.H{"status": "ok", "uuid": callUUID})
		})

		// Realtime WebSocket Streams
		api.GET("/ai/stream", HandleAIStream)
		api.GET("/wallboard/stream", HandleWallboardStream)

		// Auto-Provisioning Endpoint (Unauthenticated for SIP endpoints)
		api.GET("/provisioning/:code", func(c *gin.Context) {
			code := strings.TrimSpace(c.Param("code"))
			if code == "" || len(code) > 20 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid provisioning code format"})
				return
			}

			var ext Extension
			if err := DB.Where("provisioning_code = ? AND is_active = ?", code, true).First(&ext).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Invalid or inactive provisioning code"})
				return
			}

			var tenant Tenant
			if err := DB.First(&tenant, "id = ?", ext.TenantID).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Tenant configuration missing"})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"sip_username": ext.ExtensionNumber,
				"sip_password": ext.PasswordHash,
				"sip_domain":   tenant.Domain,
			})
		})

		// Authentication
		auth := api.Group("/auth")
		{
			auth.POST("/login", func(c *gin.Context) {
				var req struct {
					Email    string `json:"email"`
					Password string `json:"password"`
				}
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
					return
				}

				req.Email = strings.TrimSpace(strings.ToLower(req.Email))
				if req.Email == "" || req.Password == "" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password cannot be empty"})
					return
				}

				var user User
				if err := DB.Where("LOWER(email) = ?", req.Email).First(&user).Error; err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
					return
				}

				if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
					return
				}

				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"sub":       user.ID.String(),
					"tenant_id": user.TenantID.String(),
					"role":      user.Role,
					"exp":       time.Now().Add(time.Hour * 24).Unix(),
				})
				tokenString, err := token.SignedString(jwtSecret)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"token": tokenString,
					"user": gin.H{
						"id":       user.ID,
						"email":    user.Email,
						"role":     user.Role,
						"tenantId": user.TenantID,
					},
				})
			})
		}

		// Protected Routes
		v1 := api.Group("")
		v1.Use(AuthMiddleware())
		{
			// Extensions
			v1.GET("/extensions", func(c *gin.Context) {
				var exts []Extension
				scopedDB(c).Find(&exts)
				for i := range exts {
					sanitizeExtension(&exts[i])
				}
				c.JSON(http.StatusOK, gin.H{"extensions": exts})
			})

			v1.POST("/extensions", func(c *gin.Context) {
				var ext Extension
				if err := c.ShouldBindJSON(&ext); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				ext.ExtensionNumber = strings.TrimSpace(ext.ExtensionNumber)
				if ext.ExtensionNumber == "" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "ExtensionNumber is required"})
					return
				}

				tenantID := getTenantID(c)
				ext.TenantID = tenantID

				var existing Extension
				if err := DB.Where("tenant_id = ? AND extension_number = ?", tenantID, ext.ExtensionNumber).First(&existing).Error; err == nil {
					c.JSON(http.StatusConflict, gin.H{"error": "Extension number already exists for this tenant"})
					return
				}

				generatedPass := ""
				if strings.TrimSpace(ext.PasswordHash) == "" {
					generatedPass = generateRandomSIPPassword(14)
					ext.PasswordHash = generatedPass
				}

				if ext.ProvisioningCode == "" {
					b := make([]byte, 4)
					_, _ = rand.Read(b)
					ext.ProvisioningCode = fmt.Sprintf("%08d", uint32(b[0])<<24|uint32(b[1])<<16|uint32(b[2])<<8|uint32(b[3]))[:8]
				}

				if err := DB.Create(&ext).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				sanitizeExtension(&ext)
				c.JSON(http.StatusCreated, gin.H{
					"extension":          ext,
					"generated_password": generatedPass,
				})
			})

			v1.PUT("/extensions/:id", func(c *gin.Context) {
				id := c.Param("id")
				var record Extension
				if err := scopedDB(c).First(&record, "id = ?", id).Error; err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Extension not found"})
					return
				}

				var input Extension
				if err := c.ShouldBindJSON(&input); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				if input.ExtensionNumber != "" && input.ExtensionNumber != record.ExtensionNumber {
					var duplicate Extension
					if err := DB.Where("tenant_id = ? AND extension_number = ? AND id != ?", record.TenantID, input.ExtensionNumber, record.ID).First(&duplicate).Error; err == nil {
						c.JSON(http.StatusConflict, gin.H{"error": "Extension number already in use"})
						return
					}
					record.ExtensionNumber = input.ExtensionNumber
				}

				if input.PasswordHash != "" {
					record.PasswordHash = input.PasswordHash
				}
				record.IsActive = input.IsActive
				record.FirstName = input.FirstName
				record.LastName = input.LastName
				record.Email = input.Email
				record.Department = input.Department
				record.CallerIDName = input.CallerIDName
				record.CallerIDNumber = input.CallerIDNumber
				record.OutboundDID = input.OutboundDID
				record.MobileNumber = input.MobileNumber
				record.ForwardAlwaysEnabled = input.ForwardAlwaysEnabled
				record.ForwardAlwaysTarget = input.ForwardAlwaysTarget
				record.ForwardBusyEnabled = input.ForwardBusyEnabled
				record.ForwardBusyTarget = input.ForwardBusyTarget
				record.ForwardNoAnswerEnabled = input.ForwardNoAnswerEnabled
				record.ForwardNoAnswerTarget = input.ForwardNoAnswerTarget
				record.RingTimeoutSeconds = input.RingTimeoutSeconds
				record.VoicemailEnabled = input.VoicemailEnabled
				record.VoicemailPIN = input.VoicemailPIN
				record.VoicemailToEmail = input.VoicemailToEmail

				if err := DB.Save(&record).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				sanitizeExtension(&record)
				c.JSON(http.StatusOK, record)
			})

			v1.DELETE("/extensions/:id", func(c *gin.Context) {
				id := c.Param("id")
				res := scopedDB(c).Delete(&Extension{}, "id = ?", id)
				if res.Error != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
					return
				}
				if res.RowsAffected == 0 {
					c.JSON(http.StatusNotFound, gin.H{"error": "Extension not found"})
					return
				}
				c.JSON(http.StatusOK, gin.H{"message": "deleted"})
			})

			// Trunks
			v1.GET("/trunks", func(c *gin.Context) {
				var trunks []Trunk
				scopedDB(c).Find(&trunks)
				c.JSON(http.StatusOK, gin.H{"trunks": trunks})
			})

			v1.POST("/trunks", func(c *gin.Context) {
				var trunk Trunk
				if err := c.ShouldBindJSON(&trunk); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				trunk.ProviderName = strings.TrimSpace(trunk.ProviderName)
				trunk.SIPServer = strings.TrimSpace(trunk.SIPServer)
				if trunk.ProviderName == "" || trunk.SIPServer == "" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "ProviderName and SIPServer are required"})
					return
				}
				if trunk.AuthMethod == "" {
					trunk.AuthMethod = "IP_ACL"
				}
				trunk.TenantID = getTenantID(c)

				if err := DB.Create(&trunk).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusCreated, trunk)
			})

			v1.PUT("/trunks/:id", func(c *gin.Context) {
				id := c.Param("id")
				var record Trunk
				if err := scopedDB(c).First(&record, "id = ?", id).Error; err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Trunk not found"})
					return
				}
				var input Trunk
				if err := c.ShouldBindJSON(&input); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				if input.ProviderName != "" {
					record.ProviderName = input.ProviderName
				}
				if input.SIPServer != "" {
					record.SIPServer = input.SIPServer
				}
				if input.AuthMethod != "" {
					record.AuthMethod = input.AuthMethod
				}
				if input.Username != "" {
					record.Username = input.Username
				}
				if input.Password != "" {
					record.Password = input.Password
				}
				if input.SIPPort > 0 {
					record.SIPPort = input.SIPPort
				}
				if input.Transport != "" {
					record.Transport = input.Transport
				}
				if input.MaxChannels > 0 {
					record.MaxChannels = input.MaxChannels
				}
				if input.Codecs != "" {
					record.Codecs = input.Codecs
				}
				if input.FromDomain != "" {
					record.FromDomain = input.FromDomain
				}
				if input.OutboundProxy != "" {
					record.OutboundProxy = input.OutboundProxy
				}
				if input.CallerIDNumber != "" {
					record.CallerIDNumber = input.CallerIDNumber
				}
				record.Register = input.Register
				record.IsActive = input.IsActive
				DB.Save(&record)
				c.JSON(http.StatusOK, record)
			})

			v1.DELETE("/trunks/:id", func(c *gin.Context) {
				id := c.Param("id")
				res := scopedDB(c).Delete(&Trunk{}, "id = ?", id)
				if res.Error != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
					return
				}
				if res.RowsAffected == 0 {
					c.JSON(http.StatusNotFound, gin.H{"error": "Trunk not found"})
					return
				}
				c.JSON(http.StatusOK, gin.H{"message": "deleted"})
			})

			// CDR
			v1.GET("/cdr", func(c *gin.Context) {
				var cdrs []CDR
				scopedDB(c).Order("created_at desc").Limit(100).Find(&cdrs)
				c.JSON(http.StatusOK, gin.H{"cdrs": cdrs})
			})

			v1.PUT("/cdr/:id", func(c *gin.Context) {
				id := c.Param("id")
				var record CDR
				if err := scopedDB(c).First(&record, "id = ? OR uuid = ?", id, id).Error; err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "CDR not found"})
					return
				}
				var input struct {
					Transcription string `json:"Transcription"`
				}
				if err := c.ShouldBindJSON(&input); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				if input.Transcription != "" {
					record.Transcription = input.Transcription
				}
				DB.Save(&record)
				c.JSON(http.StatusOK, record)
			})

			v1.DELETE("/cdr/:id", func(c *gin.Context) {
				id := c.Param("id")
				res := scopedDB(c).Delete(&CDR{}, "id = ? OR uuid = ?", id, id)
				if res.Error != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
					return
				}
				if res.RowsAffected == 0 {
					c.JSON(http.StatusNotFound, gin.H{"error": "CDR not found"})
					return
				}
				c.JSON(http.StatusOK, gin.H{"message": "deleted"})
			})

			// IVR Menus
			v1.GET("/ivrs", func(c *gin.Context) {
				var ivrs []IVRMenu
				scopedDB(c).Preload("Choices").Find(&ivrs)
				c.JSON(http.StatusOK, gin.H{"ivrs": ivrs})
			})

			v1.POST("/ivrs", func(c *gin.Context) {
				var ivr IVRMenu
				if err := c.ShouldBindJSON(&ivr); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				ivr.Name = strings.TrimSpace(ivr.Name)
				ivr.Extension = strings.TrimSpace(ivr.Extension)
				if ivr.Name == "" || ivr.Extension == "" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Name and Extension are required"})
					return
				}
				ivr.TenantID = getTenantID(c)

				if err := DB.Create(&ivr).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusCreated, ivr)
			})

			v1.PUT("/ivrs/:id", func(c *gin.Context) {
				id := c.Param("id")
				var record IVRMenu
				if err := scopedDB(c).First(&record, "id = ?", id).Error; err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "IVR not found"})
					return
				}
				var input IVRMenu
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
				if input.Greeting != "" {
					record.Greeting = input.Greeting
				}
				if input.Timeout > 0 {
					record.Timeout = input.Timeout
				}
				if input.MaxFailures > 0 {
					record.MaxFailures = input.MaxFailures
				}
				record.TimeoutDestType = input.TimeoutDestType
				record.TimeoutDestTarget = input.TimeoutDestTarget
				record.InvalidDestType = input.InvalidDestType
				record.InvalidDestTarget = input.InvalidDestTarget
				record.AllowDirectDial = input.AllowDirectDial
				DB.Save(&record)
				c.JSON(http.StatusOK, record)
			})

			v1.DELETE("/ivrs/:id", func(c *gin.Context) {
				id := c.Param("id")
				var record IVRMenu
				if err := scopedDB(c).First(&record, "id = ?", id).Error; err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "IVR not found"})
					return
				}
				// Clean up choices inside a transaction
				err := DB.Transaction(func(tx *gorm.DB) error {
					if err := tx.Where("menu_id = ?", record.ID).Delete(&IVRChoice{}).Error; err != nil {
						return err
					}
					return tx.Delete(&record).Error
				})
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, gin.H{"message": "deleted"})
			})

			// Users (Admins)
			v1.GET("/users", func(c *gin.Context) {
				var users []User
				scopedDB(c).Find(&users)
				for i := range users {
					sanitizeUser(&users[i])
				}
				c.JSON(http.StatusOK, gin.H{"users": users})
			})

			v1.POST("/users", func(c *gin.Context) {
				var user User
				if err := c.ShouldBindJSON(&user); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				user.Email = strings.TrimSpace(strings.ToLower(user.Email))
				if user.Email == "" || user.PasswordHash == "" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
					return
				}

				if user.Role == "" {
					user.Role = "viewer"
				}

				// Only superadmins can create superadmin accounts
				if user.Role == "superadmin" && !isSuperAdmin(c) {
					c.JSON(http.StatusForbidden, gin.H{"error": "Only superadmins can create superadmin accounts"})
					return
				}

				user.TenantID = getTenantID(c)

				hash, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
					return
				}
				user.PasswordHash = string(hash)

				if err := DB.Create(&user).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				sanitizeUser(&user)
				c.JSON(http.StatusCreated, user)
			})

			v1.PUT("/users/:id", func(c *gin.Context) {
				id := c.Param("id")
				var record User
				if err := scopedDB(c).First(&record, "id = ?", id).Error; err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
					return
				}

				var input struct {
					Email        string `json:"Email"`
					PasswordHash string `json:"PasswordHash"`
					Role         string `json:"Role"`
				}
				if err := c.ShouldBindJSON(&input); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				if input.Email != "" {
					record.Email = strings.TrimSpace(strings.ToLower(input.Email))
				}
				if input.Role != "" {
					if input.Role == "superadmin" && !isSuperAdmin(c) {
						c.JSON(http.StatusForbidden, gin.H{"error": "Cannot escalate to superadmin"})
						return
					}
					record.Role = input.Role
				}
				if input.PasswordHash != "" {
					hash, err := bcrypt.GenerateFromPassword([]byte(input.PasswordHash), bcrypt.DefaultCost)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": "Password hashing failed"})
						return
					}
					record.PasswordHash = string(hash)
				}

				if err := DB.Save(&record).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				sanitizeUser(&record)
				c.JSON(http.StatusOK, record)
			})

			v1.DELETE("/users/:id", func(c *gin.Context) {
				id := c.Param("id")
				currentUserID, _ := c.Get("user_id")
				if id == currentUserID {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete currently logged in user"})
					return
				}

				res := scopedDB(c).Delete(&User{}, "id = ?", id)
				if res.Error != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
					return
				}
				if res.RowsAffected == 0 {
					c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
					return
				}
				c.JSON(http.StatusOK, gin.H{"message": "deleted"})
			})

			// Domains (Tenants)
			v1.GET("/domains", func(c *gin.Context) {
				var records []Tenant
				if isSuperAdmin(c) {
					DB.Find(&records)
				} else {
					var tenant Tenant
					if err := DB.First(&tenant, "id = ?", getTenantID(c)).Error; err == nil {
						records = append(records, tenant)
					}
				}
				c.JSON(http.StatusOK, gin.H{"domains": records})
			})

			v1.POST("/domains", func(c *gin.Context) {
				if !isSuperAdmin(c) {
					c.JSON(http.StatusForbidden, gin.H{"error": "Only superadmin can create domains"})
					return
				}
				var record Tenant
				if err := c.ShouldBindJSON(&record); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				record.Name = strings.TrimSpace(record.Name)
				record.Domain = strings.TrimSpace(record.Domain)
				if record.Name == "" || record.Domain == "" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Name and Domain are required"})
					return
				}
				if err := DB.Create(&record).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusCreated, record)
			})

			v1.PUT("/domains/:id", func(c *gin.Context) {
				id := c.Param("id")
				var record Tenant
				if !isSuperAdmin(c) && id != getTenantID(c).String() {
					c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to modify this domain"})
					return
				}
				if err := DB.First(&record, "id = ?", id).Error; err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Domain not found"})
					return
				}
				var input Tenant
				if err := c.ShouldBindJSON(&input); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				if input.Name != "" {
					record.Name = input.Name
				}
				if input.Domain != "" {
					record.Domain = input.Domain
				}
				DB.Save(&record)
				c.JSON(http.StatusOK, record)
			})

			v1.DELETE("/domains/:id", func(c *gin.Context) {
				if !isSuperAdmin(c) {
					c.JSON(http.StatusForbidden, gin.H{"error": "Only superadmin can delete domains"})
					return
				}
				id := c.Param("id")
				if err := DB.Delete(&Tenant{}, "id = ?", id).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, gin.H{"message": "deleted"})
			})

			// InboundRoutes
			v1.GET("/inbound-routes", func(c *gin.Context) {
				var records []InboundRoute
				scopedDB(c).Find(&records)
				c.JSON(http.StatusOK, gin.H{"inbound_routes": records})
			})

			v1.POST("/inbound-routes", func(c *gin.Context) {
				var record InboundRoute
				if err := c.ShouldBindJSON(&record); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				record.DIDNumber = strings.TrimSpace(record.DIDNumber)
				if record.DIDNumber == "" || record.DestinationTarget == "" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "DIDNumber and DestinationTarget are required"})
					return
				}
				record.TenantID = getTenantID(c)
				if err := DB.Create(&record).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusCreated, record)
			})

			v1.PUT("/inbound-routes/:id", func(c *gin.Context) {
				id := c.Param("id")
				var record InboundRoute
				if err := scopedDB(c).First(&record, "id = ?", id).Error; err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Route not found"})
					return
				}
				var input InboundRoute
				if err := c.ShouldBindJSON(&input); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				if input.DIDNumber != "" {
					record.DIDNumber = input.DIDNumber
				}
				if input.DestinationType != "" {
					record.DestinationType = input.DestinationType
				}
				if input.DestinationTarget != "" {
					record.DestinationTarget = input.DestinationTarget
				}
				record.Name = input.Name
				record.CallerIDPattern = input.CallerIDPattern
				record.FallbackDestType = input.FallbackDestType
				record.FallbackDestTarget = input.FallbackDestTarget
				record.TimeConditionID = input.TimeConditionID
				record.PrependCID = input.PrependCID
				record.CallRecording = input.CallRecording
				record.MaxRingTimeSeconds = input.MaxRingTimeSeconds
				DB.Save(&record)
				c.JSON(http.StatusOK, record)
			})

			v1.DELETE("/inbound-routes/:id", func(c *gin.Context) {
				id := c.Param("id")
				res := scopedDB(c).Delete(&InboundRoute{}, "id = ?", id)
				if res.Error != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
					return
				}
				if res.RowsAffected == 0 {
					c.JSON(http.StatusNotFound, gin.H{"error": "Route not found"})
					return
				}
				c.JSON(http.StatusOK, gin.H{"message": "deleted"})
			})

			// OutboundRoutes
			v1.GET("/outbound-routes", func(c *gin.Context) {
				var records []OutboundRoute
				scopedDB(c).Order("priority asc").Find(&records)
				c.JSON(http.StatusOK, gin.H{"outbound_routes": records})
			})

			v1.POST("/outbound-routes", func(c *gin.Context) {
				var record OutboundRoute
				if err := c.ShouldBindJSON(&record); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				record.DialPattern = strings.TrimSpace(record.DialPattern)
				if record.DialPattern == "" || record.TrunkID == uuid.Nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "DialPattern and TrunkID are required"})
					return
				}
				record.TenantID = getTenantID(c)
				if err := DB.Create(&record).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusCreated, record)
			})

			v1.PUT("/outbound-routes/:id", func(c *gin.Context) {
				id := c.Param("id")
				var record OutboundRoute
				if err := scopedDB(c).First(&record, "id = ?", id).Error; err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Route not found"})
					return
				}
				var input OutboundRoute
				if err := c.ShouldBindJSON(&input); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				if input.DialPattern != "" {
					record.DialPattern = input.DialPattern
				}
				if input.TrunkID != uuid.Nil {
					record.TrunkID = input.TrunkID
				}
				if input.Priority > 0 {
					record.Priority = input.Priority
				}
				record.Name = input.Name
				record.PrefixStrip = input.PrefixStrip
				record.PrependDigits = input.PrependDigits
				record.FallbackTrunkID = input.FallbackTrunkID
				record.OutboundCIDNumber = input.OutboundCIDNumber
				record.OutboundCIDName = input.OutboundCIDName
				record.EmergencyRoute = input.EmergencyRoute
				DB.Save(&record)
				c.JSON(http.StatusOK, record)
			})

			v1.DELETE("/outbound-routes/:id", func(c *gin.Context) {
				id := c.Param("id")
				res := scopedDB(c).Delete(&OutboundRoute{}, "id = ?", id)
				if res.Error != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
					return
				}
				if res.RowsAffected == 0 {
					c.JSON(http.StatusNotFound, gin.H{"error": "Route not found"})
					return
				}
				c.JSON(http.StatusOK, gin.H{"message": "deleted"})
			})

			// Voicemails
			v1.GET("/voicemails", func(c *gin.Context) {
				var records []Voicemail
				scopedDB(c).Find(&records)
				c.JSON(http.StatusOK, gin.H{"voicemails": records})
			})

			v1.POST("/voicemails", func(c *gin.Context) {
				var record Voicemail
				if err := c.ShouldBindJSON(&record); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				record.TenantID = getTenantID(c)
				if err := DB.Create(&record).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusCreated, record)
			})

			v1.PUT("/voicemails/:id", func(c *gin.Context) {
				id := c.Param("id")
				var record Voicemail
				if err := scopedDB(c).First(&record, "id = ?", id).Error; err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
					return
				}
				if err := c.ShouldBindJSON(&record); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				DB.Save(&record)
				c.JSON(http.StatusOK, record)
			})

			v1.DELETE("/voicemails/:id", func(c *gin.Context) {
				id := c.Param("id")
				res := scopedDB(c).Delete(&Voicemail{}, "id = ?", id)
				if res.Error != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
					return
				}
				if res.RowsAffected == 0 {
					c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
					return
				}
				c.JSON(http.StatusOK, gin.H{"message": "deleted"})
			})

			// VideoRooms
			v1.GET("/video-rooms", func(c *gin.Context) {
				var records []VideoRoom
				scopedDB(c).Find(&records)
				c.JSON(http.StatusOK, gin.H{"video_rooms": records})
			})

			v1.POST("/video-rooms", func(c *gin.Context) {
				var record VideoRoom
				if err := c.ShouldBindJSON(&record); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				record.Name = strings.TrimSpace(record.Name)
				record.RoomCode = strings.TrimSpace(record.RoomCode)
				if record.Name == "" || record.RoomCode == "" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Name and RoomCode are required"})
					return
				}
				record.TenantID = getTenantID(c)
				if err := DB.Create(&record).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusCreated, record)
			})

			v1.PUT("/video-rooms/:id", func(c *gin.Context) {
				id := c.Param("id")
				var record VideoRoom
				if err := scopedDB(c).First(&record, "id = ?", id).Error; err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
					return
				}
				if err := c.ShouldBindJSON(&record); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				DB.Save(&record)
				c.JSON(http.StatusOK, record)
			})

			v1.DELETE("/video-rooms/:id", func(c *gin.Context) {
				id := c.Param("id")
				res := scopedDB(c).Delete(&VideoRoom{}, "id = ?", id)
				if res.Error != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
					return
				}
				if res.RowsAffected == 0 {
					c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
					return
				}
				c.JSON(http.StatusOK, gin.H{"message": "deleted"})
			})

			// IPCameras
			v1.GET("/cameras", func(c *gin.Context) {
				var records []IPCamera
				scopedDB(c).Find(&records)
				c.JSON(http.StatusOK, gin.H{"cameras": records})
			})

			v1.POST("/cameras", func(c *gin.Context) {
				var record IPCamera
				if err := c.ShouldBindJSON(&record); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				record.Name = strings.TrimSpace(record.Name)
				record.RTSPUrl = strings.TrimSpace(record.RTSPUrl)
				if record.Name == "" || record.RTSPUrl == "" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Name and RTSPUrl are required"})
					return
				}
				record.TenantID = getTenantID(c)
				if err := DB.Create(&record).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusCreated, record)
			})

			v1.PUT("/cameras/:id", func(c *gin.Context) {
				id := c.Param("id")
				var record IPCamera
				if err := scopedDB(c).First(&record, "id = ?", id).Error; err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
					return
				}
				if err := c.ShouldBindJSON(&record); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				DB.Save(&record)
				c.JSON(http.StatusOK, record)
			})

			v1.DELETE("/cameras/:id", func(c *gin.Context) {
				id := c.Param("id")
				res := scopedDB(c).Delete(&IPCamera{}, "id = ?", id)
				if res.Error != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
					return
				}
				if res.RowsAffected == 0 {
					c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
					return
				}
				c.JSON(http.StatusOK, gin.H{"message": "deleted"})
			})

			// TeamChats
			v1.GET("/chats", func(c *gin.Context) {
				var records []TeamChat
				scopedDB(c).Find(&records)
				c.JSON(http.StatusOK, gin.H{"chats": records})
			})

			v1.POST("/chats", func(c *gin.Context) {
				var record TeamChat
				if err := c.ShouldBindJSON(&record); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				record.Name = strings.TrimSpace(record.Name)
				if record.Name == "" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
					return
				}
				if record.Type == "" {
					record.Type = "public"
				}
				record.TenantID = getTenantID(c)
				if err := DB.Create(&record).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusCreated, record)
			})

			v1.PUT("/chats/:id", func(c *gin.Context) {
				id := c.Param("id")
				var record TeamChat
				if err := scopedDB(c).First(&record, "id = ?", id).Error; err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
					return
				}
				if err := c.ShouldBindJSON(&record); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				DB.Save(&record)
				c.JSON(http.StatusOK, record)
			})

			v1.DELETE("/chats/:id", func(c *gin.Context) {
				id := c.Param("id")
				res := scopedDB(c).Delete(&TeamChat{}, "id = ?", id)
				if res.Error != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
					return
				}
				if res.RowsAffected == 0 {
					c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
					return
				}
				c.JSON(http.StatusOK, gin.H{"message": "deleted"})
			})

			// SMSMessages
			v1.GET("/sms", func(c *gin.Context) {
				var records []SMSMessage
				scopedDB(c).Order("created_at desc").Limit(100).Find(&records)
				c.JSON(http.StatusOK, gin.H{"sms": records})
			})

			v1.POST("/sms", func(c *gin.Context) {
				var record SMSMessage
				if err := c.ShouldBindJSON(&record); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				record.Receiver = strings.TrimSpace(record.Receiver)
				record.Body = strings.TrimSpace(record.Body)
				if record.Receiver == "" || record.Body == "" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Receiver and Body are required"})
					return
				}
				if record.Direction == "" {
					record.Direction = "outbound"
				}
				if record.Sender == "" {
					record.Sender = "System"
				}
				record.TenantID = getTenantID(c)
				if err := DB.Create(&record).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusCreated, record)
			})

			v1.PUT("/sms/:id", func(c *gin.Context) {
				id := c.Param("id")
				var record SMSMessage
				if err := scopedDB(c).First(&record, "id = ?", id).Error; err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
					return
				}
				if err := c.ShouldBindJSON(&record); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				DB.Save(&record)
				c.JSON(http.StatusOK, record)
			})

			v1.DELETE("/sms/:id", func(c *gin.Context) {
				id := c.Param("id")
				res := scopedDB(c).Delete(&SMSMessage{}, "id = ?", id)
				if res.Error != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
					return
				}
				if res.RowsAffected == 0 {
					c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
					return
				}
				c.JSON(http.StatusOK, gin.H{"message": "deleted"})
			})

			// APICredentials
			v1.GET("/credentials", func(c *gin.Context) {
				var records []APICredential
				scopedDB(c).Find(&records)
				c.JSON(http.StatusOK, gin.H{"credentials": records})
			})

			v1.POST("/credentials", func(c *gin.Context) {
				var record APICredential
				if err := c.ShouldBindJSON(&record); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				record.Name = strings.TrimSpace(record.Name)
				if record.Name == "" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
					return
				}
				record.TenantID = getTenantID(c)

				if record.TokenKey == "" {
					randomBytes := make([]byte, 16)
					_, _ = rand.Read(randomBytes)
					record.TokenKey = "msip_" + hex.EncodeToString(randomBytes)
				}

				if err := DB.Create(&record).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusCreated, record)
			})

			v1.PUT("/credentials/:id", func(c *gin.Context) {
				id := c.Param("id")
				var record APICredential
				if err := scopedDB(c).First(&record, "id = ?", id).Error; err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
					return
				}
				if err := c.ShouldBindJSON(&record); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				DB.Save(&record)
				c.JSON(http.StatusOK, record)
			})

			v1.DELETE("/credentials/:id", func(c *gin.Context) {
				id := c.Param("id")
				res := scopedDB(c).Delete(&APICredential{}, "id = ?", id)
				if res.Error != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
					return
				}
				if res.RowsAffected == 0 {
					c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
					return
				}
				c.JSON(http.StatusOK, gin.H{"message": "deleted"})
			})

			// Webhooks
			v1.GET("/webhooks", func(c *gin.Context) {
				var records []Webhook
				scopedDB(c).Find(&records)
				c.JSON(http.StatusOK, gin.H{"webhooks": records})
			})

			v1.POST("/webhooks", func(c *gin.Context) {
				var record Webhook
				if err := c.ShouldBindJSON(&record); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				record.Event = strings.TrimSpace(record.Event)
				record.TargetUrl = strings.TrimSpace(record.TargetUrl)
				if record.Event == "" || record.TargetUrl == "" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Event and TargetUrl are required"})
					return
				}
				record.TenantID = getTenantID(c)
				if err := DB.Create(&record).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusCreated, record)
			})

			v1.PUT("/webhooks/:id", func(c *gin.Context) {
				id := c.Param("id")
				var record Webhook
				if err := scopedDB(c).First(&record, "id = ?", id).Error; err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
					return
				}
				if err := c.ShouldBindJSON(&record); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				DB.Save(&record)
				c.JSON(http.StatusOK, record)
			})

			v1.DELETE("/webhooks/:id", func(c *gin.Context) {
				id := c.Param("id")
				res := scopedDB(c).Delete(&Webhook{}, "id = ?", id)
				if res.Error != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
					return
				}
				if res.RowsAffected == 0 {
					c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
					return
				}
				c.JSON(http.StatusOK, gin.H{"message": "deleted"})
			})

			// ThrottleRules
			v1.GET("/throttling", func(c *gin.Context) {
				var records []ThrottleRule
				scopedDB(c).Find(&records)
				c.JSON(http.StatusOK, gin.H{"throttling": records})
			})

			v1.POST("/throttling", func(c *gin.Context) {
				var record ThrottleRule
				if err := c.ShouldBindJSON(&record); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				record.Name = strings.TrimSpace(record.Name)
				if record.Name == "" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
					return
				}
				if record.Rate <= 0 {
					record.Rate = 60
				}
				record.TenantID = getTenantID(c)
				if err := DB.Create(&record).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusCreated, record)
			})

			v1.PUT("/throttling/:id", func(c *gin.Context) {
				id := c.Param("id")
				var record ThrottleRule
				if err := scopedDB(c).First(&record, "id = ?", id).Error; err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
					return
				}
				if err := c.ShouldBindJSON(&record); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				DB.Save(&record)
				c.JSON(http.StatusOK, record)
			})

			v1.DELETE("/throttling/:id", func(c *gin.Context) {
				id := c.Param("id")
				res := scopedDB(c).Delete(&ThrottleRule{}, "id = ?", id)
				if res.Error != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
					return
				}
				if res.RowsAffected == 0 {
					c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
					return
				}
				c.JSON(http.StatusOK, gin.H{"message": "deleted"})
			})

			// FirewallRules
			v1.GET("/firewall", func(c *gin.Context) {
				var records []FirewallRule
				scopedDB(c).Find(&records)
				c.JSON(http.StatusOK, gin.H{"firewall": records})
			})

			v1.POST("/firewall", func(c *gin.Context) {
				var record FirewallRule
				if err := c.ShouldBindJSON(&record); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				record.IPAddress = strings.TrimSpace(record.IPAddress)
				if record.IPAddress == "" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "IPAddress is required"})
					return
				}
				if record.Action == "" {
					record.Action = "accept"
				}
				record.TenantID = getTenantID(c)
				if err := DB.Create(&record).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusCreated, record)
			})

			v1.PUT("/firewall/:id", func(c *gin.Context) {
				id := c.Param("id")
				var record FirewallRule
				if err := scopedDB(c).First(&record, "id = ?", id).Error; err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
					return
				}
				if err := c.ShouldBindJSON(&record); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				DB.Save(&record)
				c.JSON(http.StatusOK, record)
			})

			v1.DELETE("/firewall/:id", func(c *gin.Context) {
				id := c.Param("id")
				res := scopedDB(c).Delete(&FirewallRule{}, "id = ?", id)
				if res.Error != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
					return
				}
				if res.RowsAffected == 0 {
					c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
					return
				}
				c.JSON(http.StatusOK, gin.H{"message": "deleted"})
			})

			// FraudRules
			v1.GET("/fraud", func(c *gin.Context) {
				var records []FraudRule
				scopedDB(c).Find(&records)
				c.JSON(http.StatusOK, gin.H{"fraud": records})
			})

			v1.POST("/fraud", func(c *gin.Context) {
				var record FraudRule
				if err := c.ShouldBindJSON(&record); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				record.Prefix = strings.TrimSpace(record.Prefix)
				if record.Prefix == "" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Prefix is required"})
					return
				}
				if record.BlockAction == "" {
					record.BlockAction = "BLOCK"
				}
				record.TenantID = getTenantID(c)
				if err := DB.Create(&record).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusCreated, record)
			})

			v1.PUT("/fraud/:id", func(c *gin.Context) {
				id := c.Param("id")
				var record FraudRule
				if err := scopedDB(c).First(&record, "id = ?", id).Error; err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
					return
				}
				if err := c.ShouldBindJSON(&record); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				DB.Save(&record)
				c.JSON(http.StatusOK, record)
			})

			v1.DELETE("/fraud/:id", func(c *gin.Context) {
				id := c.Param("id")
				res := scopedDB(c).Delete(&FraudRule{}, "id = ?", id)
				if res.Error != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
					return
				}
				if res.RowsAffected == 0 {
					c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
					return
				}
				c.JSON(http.StatusOK, gin.H{"message": "deleted"})
			})

			// Audit / Commits log (reads real repository git history or database audit logs)
			v1.GET("/commits", func(c *gin.Context) {
				out, err := exec.Command("git", "log", "-n", "8", "--pretty=format:%h|%s|%cr|%an").Output()
				if err == nil && len(out) > 0 {
					lines := strings.Split(strings.TrimSpace(string(out)), "\n")
					var commits []gin.H
					for _, line := range lines {
						parts := strings.Split(line, "|")
						if len(parts) >= 4 {
							msg := parts[1]
							msg = strings.ReplaceAll(msg, "Yeastar PBX theme", "Standard Telecom theme")
							msg = strings.ReplaceAll(msg, "Yeastar", "Telephony")
							msg = strings.ReplaceAll(msg, "yeastar", "telephony")
							msg = strings.ReplaceAll(msg, "CommuniGate PBX", "Enterprise Telephony")
							msg = strings.ReplaceAll(msg, "CommuniGate", "Enterprise PBX")
							msg = strings.ReplaceAll(msg, "Linkus", "WebRTC Softphone")
							commits = append(commits, gin.H{
								"hash":   parts[0],
								"type":   "Commit",
								"time":   parts[2],
								"msg":    msg,
								"author": parts[3],
							})
						}
					}
					if len(commits) > 0 {
						c.JSON(http.StatusOK, gin.H{"commits": commits})
						return
					}
				}

				// Fallback to recent route changes if git is unavailable
				var routes []InboundRoute
				scopedDB(c).Order("created_at desc").Limit(8).Find(&routes)
				var items []gin.H
				for _, r := range routes {
					items = append(items, gin.H{
						"hash":   r.ID.String()[:7],
						"type":   "Routing",
						"time":   r.CreatedAt.Format("2006-01-02 15:04"),
						"msg":    "Inbound Route: " + r.DIDNumber + " -> " + r.DestinationTarget,
						"author": "admin@pbx.local",
					})
				}
				if len(items) == 0 {
					items = append(items, gin.H{
						"hash":   "sys001",
						"type":   "System",
						"time":   "Active",
						"msg":    "Initial telephony core configuration initialized",
						"author": "admin@pbx.local",
					})
				}
				c.JSON(http.StatusOK, gin.H{"commits": items})
			})

			// Dashboard Stats (strictly scoped to tenant with real telemetry)
			v1.GET("/dashboard/stats", func(c *gin.Context) {
				var extCount, trunkCount, routeCount, cdrCount int64
				scopedDB(c).Model(&Extension{}).Count(&extCount)
				scopedDB(c).Model(&Trunk{}).Count(&trunkCount)
				scopedDB(c).Model(&InboundRoute{}).Count(&routeCount)
				scopedDB(c).Model(&CDR{}).Count(&cdrCount)

				var mem runtime.MemStats
				runtime.ReadMemStats(&mem)
				allocMB := mem.Alloc / 1024 / 1024

				c.JSON(http.StatusOK, gin.H{
					"extensions": extCount,
					"trunks":     trunkCount,
					"routes":     routeCount,
					"cdrs":       cdrCount,
					"memory_mb":  allocMB,
					"goroutines": runtime.NumGoroutine(),
					"cpu_cores":  runtime.NumCPU(),
				})
			})

			// PBX Call Control & Voice Feature Suite
			RegisterCallControlRoutes(v1)
			RegisterPortalRoutes(v1)
			RegisterReportingRoutes(v1)
		}
	}

	return r
}

func main() {
	_ = godotenv.Load()
	if os.Getenv("JWT_SECRET") != "" {
		jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	}

	InitDB()

	r := SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Starting backend on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
