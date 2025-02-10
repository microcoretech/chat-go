package websocket

import (
	"context"
	"encoding/json"

	"mbobrovskyi/chat-go/internal/chat/common"
	"mbobrovskyi/chat-go/internal/chat/domain"
)

func (e *EventHandler) updateMessagesStatusHandler(conn Connection, rawData []byte) error {
	var dto MessagesStatusDto

	if err := json.Unmarshal(rawData, &dto); err != nil {
		return err
	}

	if err := e.validate.Struct(common.ChatDomain, dto); err != nil {
		return err
	}

	chatID := conn.GetCurrentChat()
	if chatID == nil {
		return nil
	}

	if len(dto.MessageIDs) == 0 {
		return nil
	}

	if err := e.messageService.UpdateMessageStatus(
		context.Background(),
		dto.MessageIDs,
		domain.MessageStatus(dto.Status),
	); err != nil {
		return err
	}

	for _, baseConnection := range conn.GetConnector().GetConnections() {
		connection := baseConnection.(Connection)

		if !connection.IsCurrentChat(*chatID) && !connection.IsSubscribed(*chatID) {
			continue
		}

		if err := connection.SendEvent(UpdateMessagesStatusEventType, dto); err != nil {
			return err
		}
	}

	return nil
}
