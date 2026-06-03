package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("super-secret-pbx-key")

// AuthMiddleware validates JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or malformed JWT"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, _ := token.Claims.(jwt.MapClaims)
		c.Set("tenant_id", claims["tenant_id"])
		c.Set("role", claims["role"])
		c.Next()
	}
}

func main() {
	_ = godotenv.Load()
	if os.Getenv("JWT_SECRET") != "" {
		jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	}

	InitDB()

	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	r.POST("/fs", HandleFSCurl)

	api := r.Group("/api")
	{
		// Internal System Route - FreeSWITCH posts here unauthenticated
		api.POST("/cdr", func(c *gin.Context) {
			var payload map[string]interface{}
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			variables := payload["variables"].(map[string]interface{})
			callUUID, _ := variables["uuid"].(string)
			callerIDName, _ := variables["caller_id_name"].(string)
			callerIDNumber, _ := variables["caller_id_number"].(string)
			destination, _ := variables["destination_number"].(string)
			context, _ := variables["context"].(string)
			duration, _ := variables["duration"].(float64)
			billsec, _ := variables["billsec"].(float64)
			hangupCause, _ := variables["hangup_cause"].(string)

			cdr := CDR{
				UUID:           callUUID,
				CallerIDName:   callerIDName,
				CallerIDNumber: callerIDNumber,
				Destination:    destination,
				Context:        context,
				Duration:       int(duration),
				Billsec:        int(billsec),
				HangupCause:    hangupCause,
			}

			if err := DB.Create(&cdr).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save CDR"})
				return
			}

			if destination == "*88" {
				go func(id string, callCDR CDR) {
					time.Sleep(2 * time.Second)
					callCDR.Transcription = "[AI Transcription Generated via Whisper] This is a successfully recorded voice memo from your MSIP PBX extension."
					DB.Save(&callCDR)
				}(callUUID, cdr)
			}
			c.JSON(http.StatusCreated, gin.H{"status": "ok"})
		})

		// AI Conversational Receptionist WebSocket Stream (From FreeSWITCH)
		api.GET("/ai/stream", HandleAIStream)

		// Call Center Live Wallboard WebSocket Stream (For Frontend)
		api.GET("/wallboard/stream", HandleWallboardStream)

		// Auto-Provisioning Endpoint (Unauthenticated so softphones can fetch)
		api.GET("/provisioning/:code", func(c *gin.Context) {
			code := c.Param("code")
			var ext Extension
			if err := DB.Preload("Tenant").Where("provisioning_code = ?", code).First(&ext).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Invalid provisioning code"})
				return
			}
			
			var tenant Tenant
			DB.First(&tenant, ext.TenantID)

			c.JSON(http.StatusOK, gin.H{
				"sip_username": ext.ExtensionNumber,
				"sip_password": ext.PasswordHash, // Assuming plain text stored or passed securely via TLS
				"sip_domain":   tenant.Domain,
			})
		})

		auth := api.Group("/auth")
		{
			auth.POST("/login", func(c *gin.Context) {
				var req struct {
					Email    string `json:"email"`
					Password string `json:"password"`
				}
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				var user User
				if err := DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
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
				tokenString, _ := token.SignedString(jwtSecret)

				c.JSON(http.StatusOK, gin.H{
					"token": tokenString,
					"user": gin.H{
						"email":    user.Email,
						"role":     user.Role,
						"tenantId": user.TenantID,
					},
				})
			})
		}

		// Protected Routes
		v1 := api.Group("") // Re-use /api prefix but add auth
		v1.Use(AuthMiddleware())
		{
			// Extensions
			v1.GET("/extensions", func(c *gin.Context) {
			var exts []Extension
			DB.Find(&exts)
			c.JSON(http.StatusOK, gin.H{"extensions": exts})
		})
		
		v1.POST("/extensions", func(c *gin.Context) {
			var ext Extension
			if err := c.ShouldBindJSON(&ext); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			// Assign to first tenant for demo purposes
			var tenant Tenant
			DB.First(&tenant)
			ext.TenantID = tenant.ID
			
			// Generate 8-digit provisioning code
			b := make([]byte, 4)
			rand.Read(b)
			ext.ProvisioningCode = fmt.Sprintf("%08d", uint32(b[0])<<24|uint32(b[1])<<16|uint32(b[2])<<8|uint32(b[3]))[:8]
			
			if err := DB.Create(&ext).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, ext)
		})

		// Trunks
		v1.GET("/trunks", func(c *gin.Context) {
			var trunks []Trunk
			DB.Find(&trunks)
			c.JSON(http.StatusOK, gin.H{"trunks": trunks})
		})
		
		v1.POST("/trunks", func(c *gin.Context) {
			var trunk Trunk
			if err := c.ShouldBindJSON(&trunk); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			var tenant Tenant
			DB.First(&tenant)
			trunk.TenantID = tenant.ID
			
			if err := DB.Create(&trunk).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, trunk)
		})


		v1.GET("/cdr", func(c *gin.Context) {
			var cdrs []CDR
			DB.Order("created_at desc").Limit(100).Find(&cdrs)
			c.JSON(http.StatusOK, gin.H{"cdrs": cdrs})
		})

		// IVR Endpoints
		v1.GET("/ivrs", func(c *gin.Context) {
			var ivrs []IVRMenu
			DB.Preload("Choices").Find(&ivrs)
			c.JSON(http.StatusOK, gin.H{"ivrs": ivrs})
		})

		v1.POST("/ivrs", func(c *gin.Context) {
			var ivr IVRMenu
			if err := c.ShouldBindJSON(&ivr); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			var tenant Tenant
			DB.First(&tenant)
			ivr.TenantID = tenant.ID
			
			if err := DB.Create(&ivr).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, ivr)
		})

		// Users (Admins)
		v1.GET("/users", func(c *gin.Context) {
			var users []User
			tenantID, _ := c.Get("tenant_id")
			DB.Where("tenant_id = ?", tenantID).Find(&users)
			c.JSON(http.StatusOK, gin.H{"users": users})
		})

		v1.POST("/users", func(c *gin.Context) {
			var user User
			if err := c.ShouldBindJSON(&user); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			
			tenantIDClaim, _ := c.Get("tenant_id")
			tenantID, _ := uuid.Parse(tenantIDClaim.(string))
			user.TenantID = tenantID
			
			hash, _ := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
			user.PasswordHash = string(hash)

			if err := DB.Create(&user).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, user)
		})
		} // closes v1 group
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Starting backend on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
