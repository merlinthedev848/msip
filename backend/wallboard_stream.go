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
			// Generate simulated mock data
			metrics := gin.H{
				"active_calls":   rand.Intn(50) + 10,
				"agents_online":  rand.Intn(15) + 5,
				"agents_busy":    rand.Intn(10),
				"queue_length":   rand.Intn(20),
				"avg_wait_time":  rand.Intn(120), // seconds
				"sla_percentage": 85 + rand.Intn(15),
				"timestamp":      time.Now().Unix(),
			}

			if err := ws.WriteJSON(metrics); err != nil {
				log.Printf("Wallboard client disconnected: %v", err)
				return
			}
		}
	}
}
