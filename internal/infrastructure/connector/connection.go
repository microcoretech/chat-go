package connector

import "mbobrovskyi/chat-go/internal/common/domain"

type Connection interface {
	IsClosed() bool

	GetConnectionID() string

	GetConnector() Connector
	SetConnector(connector Connector)

	GetMessageChan() chan []byte
	GetCloseChan() chan struct{}

	SendEvent(eventType uint64, data any) error

	Connect()
	Close()

	GetSession() *domain.Session
}
