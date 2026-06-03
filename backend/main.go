package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"service": "pbx-backend",
		})
	})

	api := r.Group("/api")
	{
		api.GET("/extensions", func(c *gin.Context) {
			// Dummy response, will be replaced with DB logic
			c.JSON(http.StatusOK, gin.H{
				"extensions": []gin.H{
					{"id": "1", "tenant_id": "tenant-1", "ext": "101"},
					{"id": "2", "tenant_id": "tenant-1", "ext": "102"},
				},
			})
		})
	}

	// For JWT validation mocking if needed locally
	r.GET("/.well-known/jwks.json", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"keys": []interface{}{}})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Starting backend on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
