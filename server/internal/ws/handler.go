package ws

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Top level definition for all messages
// Payload structure defined in payload function files
type Message struct {
	Action  string          `json:"action"`
	Payload json.RawMessage `json:"payload"`
}

type HandlerFunc func(conn *websocket.Conn, payload json.RawMessage) error

var handlers = map[string]HandlerFunc{}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}
	defer conn.Close()

	log.Printf("client connected: %s", r.RemoteAddr)

	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}

		var msg Message
		if err := json.Unmarshal(raw, &msg); err != nil {
			log.Println("bad message:", err)
			continue
		}

		h, ok := handlers[msg.Action]
		if !ok {
			log.Printf("unknown action: %s", msg.Action)
			continue
		}

		if err := h(conn, msg.Payload); err != nil {
			log.Printf("handler error for %s: %v", msg.Action, err)
		}
	}

	log.Printf("client disconnected: %s", r.RemoteAddr)
}

func Register(action string, h HandlerFunc) {
	handlers[action] = h
}
