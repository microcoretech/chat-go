package websocket

import (
	"context"

	"mbobrovskyi/chat-go/internal/chat/domain"
	"mbobrovskyi/chat-go/internal/infrastructure/connector"
	"mbobrovskyi/chat-go/internal/infrastructure/validator"
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
