package ws

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

func init() {
	Register("room.join", handleRoomJoin)
	Register("room.leave", handleRoomLeave)
}

type roomJoinPayload struct {
	Room string `json:"room"`
}

type roomLeavePayload struct {
	Room string `json:"room"`
}

func handleRoomJoin(conn *websocket.Conn, payload json.RawMessage) error {
	var p roomJoinPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return err
	}
	log.Printf("room.join → room=%s", p.Room)
	// TODO: subscribe conn to room, spool transaction log for delivery
	return nil
}

func handleRoomLeave(conn *websocket.Conn, payload json.RawMessage) error {
	var p roomLeavePayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return err
	}
	log.Printf("room.leave → room=%s", p.Room)
	// TODO: unsubscribe conn from room
	return nil
}
