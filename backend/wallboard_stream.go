package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

// HandleWallboardStream pushes live metrics to the frontend
func HandleWallboardStream(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade websocket for Wallboard Stream: %v", err)
		return
	}
	defer ws.Close()

	log.Println("New Wallboard client connected")

	// In a real implementation, this would subscribe to a Redis PubSub channel
	// fed by FreeSWITCH ESL events (e.g. mod_callcenter agent status, queue length).
	// For now, we simulate real-time metrics pushing every 2 seconds.

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Generate simulated mock data matching the frontend expectations
			response := gin.H{
				"metrics": gin.H{
					"active_calls":   rand.Intn(50) + 10,
					"agents_online":  rand.Intn(15) + 5,
					"agents_busy":    rand.Intn(10),
					"queue_length":   rand.Intn(20),
					"avg_wait_time":  rand.Intn(120), // seconds
					"sla_percentage": 85 + rand.Intn(15),
					"timestamp":      time.Now().Unix(),
				},
				"queues": []gin.H{
					{"name": "Support Tier 1", "waiting": rand.Intn(5), "agents": 8, "sla": 85 + rand.Intn(15)},
					{"name": "Sales Global", "waiting": rand.Intn(3), "agents": 4, "sla": 90 + rand.Intn(10)},
					{"name": "Billing", "waiting": rand.Intn(2), "agents": 2, "sla": 95 + rand.Intn(5)},
				},
			}

			if err := ws.WriteJSON(response); err != nil {
				log.Printf("Wallboard client disconnected: %v", err)
				return
			}
		}
	}
}
