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

	"chat-go/internal/chat/common"
	"chat-go/internal/chat/domain"
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
