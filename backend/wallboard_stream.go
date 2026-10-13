package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wallboardUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// HandleWallboardStream pushes live call center metrics to the frontend
func HandleWallboardStream(c *gin.Context) {
	ws, err := wallboardUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade websocket for Wallboard Stream: %v", err)
		return
	}
	defer ws.Close()

	done := make(chan struct{})

	// Read loop to detect client disconnect cleanly
	go func() {
		defer close(done)
		for {
			if _, _, err := ws.ReadMessage(); err != nil {
				return
			}
		}
	}()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-c.Request.Context().Done():
			return
		case <-ticker.C:
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

			_ = ws.SetWriteDeadline(time.Now().Add(5 * time.Second))
			if err := ws.WriteJSON(response); err != nil {
				return
			}
		}
	}
}
