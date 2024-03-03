package websocket

import (
	"encoding/json"
)

type SubscribeRoomEventData struct {
	RoomID uint64 `json:"roomId"`
}

func (e *EventHandler) subscribeChatHandler(conn Connection, rawData []byte) error {
	var chatIDs []uint64

	if err := json.Unmarshal(rawData, &chatIDs); err != nil {
		return err
	}

	conn.SetSubscribedChats(chatIDs)

	return nil
}
