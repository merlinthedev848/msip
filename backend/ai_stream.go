package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var aiUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// AIStreamIncomingMsg represents incoming JSON frames from FreeSWITCH audio fork or WebRTC bridge
type AIStreamIncomingMsg struct {
	Event        string `json:"event"`         // start, speech_text, audio_frame, hangup
	UUID         string `json:"uuid"`          // Call UUID
	CallerNumber string `json:"caller_number"` // Caller ID Number
	CallerName   string `json:"caller_name"`   // Caller ID Name
	Text         string `json:"text"`          // Recognized speech text
	Transcript   string `json:"transcript"`    // Alternative transcript field
}

// AIStreamOutgoingMsg represents outgoing TTS and control frames sent to FreeSWITCH
type AIStreamOutgoingMsg struct {
	Event      string        `json:"event"`                // tts_playback, call_transfer, hangup
	Text       string        `json:"text,omitempty"`       // Synthesized text
	Voice      string        `json:"voice,omitempty"`      // Voice name/ID
	Provider   string        `json:"provider,omitempty"`   // TTS provider
	Target     string        `json:"target,omitempty"`     // Extension for call_transfer
	ToolCalls  []AIToolCall  `json:"tool_calls,omitempty"` // Any executed tools
	Action     string        `json:"action,omitempty"`
}

// HandleAIStream receives realtime audio / event frames from FreeSWITCH (via mod_audio_fork)
// or browser WebRTC simulator, processes conversational multi-turn dialogue with tool execution,
// and streams audio/actions back to the caller while persisting call session analytics.
func HandleAIStream(c *gin.Context) {
	ws, err := aiUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("[AI Stream] Failed to upgrade websocket: %v", err)
		return
	}
	defer ws.Close()

	assistantIDParam := c.Query("assistant_id")
	callerNumber := c.Query("caller_id_number")
	callUUID := c.Query("uuid")
	if callUUID == "" {
		callUUID = uuid.New().String()
	}
	if callerNumber == "" {
		callerNumber = "Unknown Caller"
	}

	// 1. Locate Assistant persona
	var assistant AIAssistant
	if assistantIDParam != "" {
		_ = DB.First(&assistant, "id = ?", assistantIDParam)
	}
	if assistant.ID == uuid.Nil {
		_ = DB.Where("is_active = ?", true).First(&assistant)
	}
	if assistant.ID == uuid.Nil {
		// Fallback default persona
		assistant = AIAssistant{
			Base:            Base{ID: uuid.New()},
			Name:            "Sarah - Virtual Assistant",
			Extension:       "890",
			ExtensionNumber: "890",
			Greeting:        "Hello, thank you for calling MSIP PBX. How may I direct your call or assist you today?",
			GreetingMessage: "Hello, thank you for calling MSIP PBX. How may I direct your call or assist you today?",
			LLMProvider:     "openai",
			LLMModel:        "gpt-4o",
			VoiceProvider:   "openai",
			TTSProvider:     "openai",
			VoiceName:       "alloy",
			VoiceID:         "alloy",
			EnableBooking:   true,
			EnableTransfer:  true,
			TransferTarget:  "1001",
		}
	}

	startTime := time.Now()
	log.Printf("[AI Stream] Active session started for call %s with assistant '%s' (Caller: %s)", callUUID, assistant.Name, callerNumber)

	// 2. Initialize Session Record in DB
	session := AICallSession{
		Base:          Base{ID: uuid.New()},
		TenantRef:     TenantRef{TenantID: assistant.TenantID},
		AssistantID:   assistant.ID,
		AssistantName: assistant.Name,
		CallUUID:      callUUID,
		CallerNumber:  callerNumber,
		ActionTaken:   "conversation",
		Sentiment:     "neutral",
		MessagesJSON:  "[]",
	}
	_ = DB.Create(&session)

	// In-memory conversation turns
	conversation := []AIMessage{}
	greeting := assistant.EffectiveGreeting()

	conversation = append(conversation, AIMessage{
		Role:    "assistant",
		Content: greeting,
		Time:    time.Now().Format("15:04:05"),
	})

	// Send initial greeting to caller
	welcomePayload := AIStreamOutgoingMsg{
		Event:    "tts_playback",
		Text:     greeting,
		Voice:    assistant.EffectiveVoice(),
		Provider: assistant.EffectiveVoiceProvider(),
	}
	welcomeBytes, _ := json.Marshal(welcomePayload)
	_ = ws.SetWriteDeadline(time.Now().Add(5 * time.Second))
	_ = ws.WriteMessage(websocket.TextMessage, welcomeBytes)

	done := make(chan struct{})

	// Handle WebSocket reading loop
	for {
		messageType, messageBytes, err := ws.ReadMessage()
		if err != nil {
			log.Printf("[AI Stream] Call %s disconnected: %v", callUUID, err)
			break
		}

		if messageType == websocket.BinaryMessage {
			// Binary PCM audio chunk received from FreeSWITCH mod_audio_fork
			// Kept active; can be routed to real-time VAD / streaming STT if configured
			continue
		}

		if messageType == websocket.TextMessage {
			var incoming AIStreamIncomingMsg
			if err := json.Unmarshal(messageBytes, &incoming); err != nil {
				// Treat plain string as speech text directly
				incoming.SpeechText(string(messageBytes))
			}

			userSpeech := incoming.Text
			if userSpeech == "" {
				userSpeech = incoming.Transcript
			}

			if incoming.CallerNumber != "" {
				callerNumber = incoming.CallerNumber
				session.CallerNumber = callerNumber
			}

			if userSpeech == "" {
				continue
			}

			log.Printf("[AI Stream] [%s] Caller said: %s", callUUID, userSpeech)

			// Append user turn
			conversation = append(conversation, AIMessage{
				Role:    "user",
				Content: userSpeech,
				Time:    time.Now().Format("15:04:05"),
			})

			// Process AI reply & execute tools (Booking, Transfer, FAQ)
			reply, tools, err := ProcessAIMessage(assistant.TenantID, assistant, conversation, callerNumber)
			if err != nil {
				reply = "I apologize, I experienced a temporary issue. Let me connect you with an operator right away."
			}

			conversation = append(conversation, AIMessage{
				Role:    "assistant",
				Content: reply,
				Time:    time.Now().Format("15:04:05"),
			})

			// Check if tool triggered booking or transfer
			var transferTarget string
			for _, tool := range tools {
				if tool.Name == "book_appointment" {
					session.ActionTaken = "booking_created"
				} else if tool.Name == "transfer_call" {
					session.ActionTaken = "transferred"
					if tgt, ok := tool.Arguments["target_extension"].(string); ok && tgt != "" {
						transferTarget = tgt
					} else if tgt, ok := tool.Arguments["destination_extension"].(string); ok && tgt != "" {
						transferTarget = tgt
					} else {
						transferTarget = assistant.TransferTarget
					}
				}
			}

			// Send assistant TTS response back to FreeSWITCH / caller
			outMsg := AIStreamOutgoingMsg{
				Event:     "tts_playback",
				Text:      reply,
				Voice:     assistant.EffectiveVoice(),
				Provider:  assistant.EffectiveVoiceProvider(),
				ToolCalls: tools,
			}
			outBytes, _ := json.Marshal(outMsg)
			_ = ws.SetWriteDeadline(time.Now().Add(5 * time.Second))
			_ = ws.WriteMessage(websocket.TextMessage, outBytes)

			// If call transfer was requested, emit transfer instruction
			if transferTarget != "" {
				transferMsg := AIStreamOutgoingMsg{
					Event:  "call_transfer",
					Target: transferTarget,
					Action: fmt.Sprintf("transfer %s XML default", transferTarget),
				}
				transferBytes, _ := json.Marshal(transferMsg)
				_ = ws.SetWriteDeadline(time.Now().Add(5 * time.Second))
				_ = ws.WriteMessage(websocket.TextMessage, transferBytes)
				log.Printf("[AI Stream] [%s] Call transferred to %s", callUUID, transferTarget)
			}
		}
	}

	close(done)

	// 3. Post-call Session Finalization & Analytics
	duration := int(time.Since(startTime).Seconds())
	session.DurationSec = duration

	var fullConversationText strings.Builder
	for _, m := range conversation {
		fullConversationText.WriteString(fmt.Sprintf("%s: %s\n", strings.ToUpper(m.Role), m.Content))
	}

	transcriptStr := fullConversationText.String()
	summary, sentiment := generateSummaryAndSentiment(transcriptStr)
	session.Summary = summary
	session.Sentiment = sentiment

	msgsBytes, _ := json.Marshal(conversation)
	session.MessagesJSON = string(msgsBytes)

	_ = DB.Save(&session)
	log.Printf("[AI Stream] Session %s finalized. Duration: %ds, Sentiment: %s, Action: %s",
		callUUID, duration, sentiment, session.ActionTaken)
}

func (m *AIStreamIncomingMsg) SpeechText(raw string) {
	clean := strings.TrimSpace(raw)
	if strings.HasPrefix(clean, "{") {
		return
	}
	m.Text = clean
}
