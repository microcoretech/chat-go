package websocket

import (
	"context"
	"encoding/json"
)

func (e *EventHandler) createMessageHandler(conn Connection, rawData []byte) error {
	dto := MessageDto{}
	if err := json.Unmarshal(rawData, &dto); err != nil {
		return err
	}

	uuid := dto.UUID

	chatID := conn.GetCurrentChat()
	if chatID == nil {
		return nil
	}

	newMessage := MessageFromCreateDto(dto)
	newMessage.ChatID = *chatID
	newMessage.CreatedBy = conn.GetSession().User.ID

	message, err := e.messageService.CreateMessage(context.Background(), newMessage)
	if err != nil {
		return err
	}

	if message == nil {
		return nil
	}

	dto = MessageToDto(*message)

	for _, baseConnection := range conn.GetConnector().GetConnections() {
		connection := baseConnection.(Connection)

		if connection.GetConnectionID() == conn.GetConnectionID() {
			continue
		}

		if !connection.IsCurrentChat(*chatID) && !connection.IsSubscribed(*chatID) {
			continue
		}

		if err := connection.SendEvent(CreateMessageEventType, dto); err != nil {
			return err
		}
	}

	dto.UUID = uuid
	if err := conn.SendEvent(CreateMessageEventType, dto); err != nil {
		return err
	}

	return nil
}
