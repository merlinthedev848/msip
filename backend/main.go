package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables if .env exists
	_ = godotenv.Load()

	// Initialize Database
	InitDB()

	r := gin.Default()

	// Configure CORS so Vite frontend can access the API
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true // For local dev, allow all
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "pbx-backend",
		})
	})

	// FreeSWITCH mod_xml_curl endpoint
	r.POST("/fs", HandleFSCurl)

	api := r.Group("/api")
	{
		// Extensions
		api.GET("/extensions", func(c *gin.Context) {
			var exts []Extension
			DB.Find(&exts)
			c.JSON(http.StatusOK, gin.H{"extensions": exts})
		})
		
		api.POST("/extensions", func(c *gin.Context) {
			var ext Extension
			if err := c.ShouldBindJSON(&ext); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			// Assign to first tenant for demo purposes
			var tenant Tenant
			DB.First(&tenant)
			ext.TenantID = tenant.ID
			
			if err := DB.Create(&ext).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, ext)
		})

		// Trunks
		api.GET("/trunks", func(c *gin.Context) {
			var trunks []Trunk
			DB.Find(&trunks)
			c.JSON(http.StatusOK, gin.H{"trunks": trunks})
		})
		
		api.POST("/trunks", func(c *gin.Context) {
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

		// CDR Ingestion from FreeSWITCH mod_json_cdr
		api.POST("/cdr", func(c *gin.Context) {
			var payload map[string]interface{}
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// Parse relevant fields from FreeSWITCH JSON CDR
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
			c.JSON(http.StatusCreated, gin.H{"status": "ok"})
		})

		api.GET("/cdr", func(c *gin.Context) {
			var cdrs []CDR
			DB.Order("created_at desc").Limit(100).Find(&cdrs)
			c.JSON(http.StatusOK, gin.H{"cdrs": cdrs})
		})

		// IVR Endpoints
		api.GET("/ivrs", func(c *gin.Context) {
			var ivrs []IVRMenu
			DB.Preload("Choices").Find(&ivrs)
			c.JSON(http.StatusOK, gin.H{"ivrs": ivrs})
		})

		api.POST("/ivrs", func(c *gin.Context) {
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
