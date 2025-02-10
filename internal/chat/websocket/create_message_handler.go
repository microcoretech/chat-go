// Copyright 2025 Mykhailo Bobrovskyi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
