package ws

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

func init() {
	Register("user.connect", handleUserConnect)
}

type userConnectPayload struct {
	UserID string `json:"user_id"`
}

func handleUserConnect(conn *websocket.Conn, payload json.RawMessage) error {
	var p userConnectPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return err
	}
	log.Printf("user.connect → user=%s", p.UserID)
	// TODO: authenticate user, stream spooled transaction log
	return nil
}
