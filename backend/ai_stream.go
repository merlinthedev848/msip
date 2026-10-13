package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var aiUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// HandleAIStream receives realtime audio from FreeSWITCH (e.g. via mod_audio_fork)
// and streams it to the LLM. It also receives TTS audio from the LLM and streams
// it back to FreeSWITCH.
func HandleAIStream(c *gin.Context) {
	ws, err := aiUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade websocket for AI Stream: %v", err)
		return
	}
	defer ws.Close()

	log.Println("New AI Receptionist audio stream connected from FreeSWITCH")

	done := make(chan struct{})

	// Send initial greeting prompt
	go func() {
		select {
		case <-time.After(1 * time.Second):
			welcomeMsg := `{"event": "tts_playback", "text": "Welcome to MSIP PBX. I am your AI assistant. How can I help you today?"}`
			_ = ws.SetWriteDeadline(time.Now().Add(5 * time.Second))
			_ = ws.WriteMessage(websocket.TextMessage, []byte(welcomeMsg))
		case <-done:
			return
		}
	}()

	for {
		messageType, message, err := ws.ReadMessage()
		if err != nil {
			close(done)
			log.Printf("AI Stream disconnected: %v", err)
			break
		}

		if messageType == websocket.BinaryMessage {
			// Raw audio payload from caller
			// Can be forwarded to external audio processing / LLM websocket
			_ = message
		} else if messageType == websocket.TextMessage {
			// JSON metadata (e.g. Call UUID, caller ID)
			log.Printf("AI Stream Metadata Received: %s", string(message))
		}
	}
}
