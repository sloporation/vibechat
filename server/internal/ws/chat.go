package ws

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

func init() {
	Register("chat.send", handleChatSend)
	Register("chat.edit", handleChatEdit)
}

type chatSendPayload struct {
	To      string `json:"to"`
	Message string `json:"message"`
}

type chatEditPayload struct {
	MessageID string `json:"message_id"`
	Message   string `json:"message"`
}

func handleChatSend(conn *websocket.Conn, payload json.RawMessage) error {
	var p chatSendPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return err
	}
	log.Printf("chat.send → to=%s msg=%s", p.To, p.Message)
	// TODO: route message to recipient or federation target
	return nil
}

func handleChatEdit(conn *websocket.Conn, payload json.RawMessage) error {
	var p chatEditPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return err
	}
	log.Printf("chat.edit → id=%s msg=%s", p.MessageID, p.Message)
	// TODO: append edit addendum to transaction log
	return nil
}
