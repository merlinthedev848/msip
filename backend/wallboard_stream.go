package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wallboardUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// HandleWallboardStream pushes live call center metrics to the frontend from live database records
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
			var activeExts int64
			DB.Model(&Extension{}).Where("is_active = ?", true).Count(&activeExts)

			var todayCalls int64
			DB.Model(&CDR{}).Where("created_at >= ?", time.Now().Add(-24*time.Hour)).Count(&todayCalls)

			var answeredCalls int64
			DB.Model(&CDR{}).Where("created_at >= ? AND disposition = ?", time.Now().Add(-24*time.Hour), "ANSWERED").Count(&answeredCalls)

			var avgDuration float64
			DB.Model(&CDR{}).Where("created_at >= ? AND disposition = ?", time.Now().Add(-24*time.Hour), "ANSWERED").Select("COALESCE(AVG(duration), 0)").Scan(&avgDuration)

			sla := 100
			if todayCalls > 0 {
				sla = int((answeredCalls * 100) / todayCalls)
			}

			// Find real Ring Groups
			var ivrs []IVRMenu
			DB.Preload("Choices").Where("greeting = ? OR name LIKE ?", "ring_group", "Ring Group:%").Find(&ivrs)

			queueStats := make([]gin.H, 0)
			for _, ivr := range ivrs {
				displayName := strings.TrimPrefix(ivr.Name, "Ring Group: ")
				if displayName == "" {
					displayName = ivr.Name
				}
				queueStats = append(queueStats, gin.H{
					"name":    fmt.Sprintf("%s (Ext %s)", displayName, ivr.Extension),
					"waiting": 0,
					"agents":  len(ivr.Choices),
					"sla":     sla,
				})
			}

			response := gin.H{
				"metrics": gin.H{
					"active_calls":   0,
					"agents_online":  activeExts,
					"agents_busy":    0,
					"queue_length":   0,
					"avg_wait_time":  int(avgDuration),
					"sla_percentage": sla,
					"timestamp":      time.Now().Unix(),
				},
				"queues": queueStats,
			}

			_ = ws.SetWriteDeadline(time.Now().Add(5 * time.Second))
			if err := ws.WriteJSON(response); err != nil {
				return
			}
		}
	}
}
