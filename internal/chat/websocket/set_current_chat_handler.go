package websocket

import (
	"encoding/json"
)

func (e *EventHandler) setCurrentChatHandler(conn Connection, rawData []byte) error {
	var chatID uint64

	if err := json.Unmarshal(rawData, &chatID); err != nil {
		return err
	}

	conn.SetCurrentChat(&chatID)

	return nil
}
