// Copyright MicroCore Tech
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

	"chat-go/internal/chat/domain"
	"chat-go/internal/infrastructure/connector"
	"chat-go/internal/infrastructure/validator"
)

type MessageService interface {
	CreateMessage(ctx context.Context, message domain.Message) (*domain.Message, error)
	UpdateMessageStatus(ctx context.Context, messageIDs []uint64, status domain.MessageStatus) error
}

type EventHandler struct {
	validate       validator.Validate
	messageService MessageService
}

func (e *EventHandler) HandleEvent(baseConn connector.Connection, event connector.Event) error {
	conn := baseConn.(Connection)

	switch event.Type {
	case SubscribeChatsEventType:
		return e.subscribeChatHandler(conn, event.Data)
	case UnsubscribeChatsEventType:
		return e.unsubscribeRoomHandler(conn, event.Data)
	case SetCurrentChatEventType:
		return e.setCurrentChatHandler(conn, event.Data)
	case UnsetCurrentChatEventType:
		return e.unsetCurrentChatHandler(conn, event.Data)
	case CreateMessageEventType:
		return e.createMessageHandler(conn, event.Data)
	case EditMessageEventType:
	case DeleteMessageEventType:
	case UpdateMessagesStatusEventType:
		return e.updateMessagesStatusHandler(conn, event.Data)
	}

	return nil
}

func NewEventHandler(
	validate validator.Validate,
	messageService MessageService,
) *EventHandler {
	return &EventHandler{
		validate:       validate,
		messageService: messageService,
	}
}
