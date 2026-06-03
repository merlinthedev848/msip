package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// HandleAIStream receives realtime audio from FreeSWITCH (e.g. via mod_audio_fork)
// and streams it to the LLM. It also receives TTS audio from the LLM and streams
// it back to FreeSWITCH.
func HandleAIStream(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade websocket for AI Stream: %v", err)
		return
	}
	defer ws.Close()

	log.Println("New AI Receptionist audio stream connected from FreeSWITCH")

	// Start a mock OpenAI LLM response loop
	go func() {
		time.Sleep(2 * time.Second)
		welcomeMsg := `{"event": "tts_playback", "text": "Welcome to MSIP PBX. I am your AI assistant. How can I help you today?"}`
		ws.WriteMessage(websocket.TextMessage, []byte(welcomeMsg))
	}()

	for {
		messageType, message, err := ws.ReadMessage()
		if err != nil {
			log.Printf("AI Stream disconnected: %v", err)
			break
		}

		if messageType == websocket.BinaryMessage {
			// Raw audio payload from caller
			// TODO: Forward this binary frame to OpenAI Realtime API WebSocket
			_ = message
		} else if messageType == websocket.TextMessage {
			// JSON metadata (e.g., Call UUID, caller ID)
			log.Printf("AI Stream Metadata Received: %s", string(message))
		}
	}
}
