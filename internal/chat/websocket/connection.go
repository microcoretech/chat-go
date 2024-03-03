package websocket

import (
	"chat/internal/common/domain"
	"chat/internal/infrastructure/connector"
	"github.com/fasthttp/websocket"
	"golang.org/x/exp/slices"
)

type Connection interface {
	connector.Connection

	GetSubscribedChats() []uint64
	SetSubscribedChats(ids []uint64)
	IsSubscribed(chatID uint64) bool

	GetCurrentChat() *uint64
	SetCurrentChat(id *uint64)
	IsCurrentChat(chatID uint64) bool
}

type connectionImpl struct {
	connector.Connection

	subscribedChats []uint64
	currentChat     *uint64
}

func (c *connectionImpl) GetSubscribedChats() []uint64 {
	return c.subscribedChats
}

func (c *connectionImpl) SetSubscribedChats(ids []uint64) {
	c.subscribedChats = ids
}

func (c *connectionImpl) IsSubscribed(chatID uint64) bool {
	return slices.Contains(c.GetSubscribedChats(), chatID)
}

func (c *connectionImpl) GetCurrentChat() *uint64 {
	return c.currentChat
}

func (c *connectionImpl) SetCurrentChat(id *uint64) {
	c.currentChat = id
}

func (c *connectionImpl) IsCurrentChat(chatID uint64) bool {
	if c.GetCurrentChat() == nil {
		return false
	}

	return *c.GetCurrentChat() == chatID
}

func NewConnection(conn *websocket.Conn, session *domain.Session) connector.Connection {
	return &connectionImpl{
		Connection: connector.NewWebSocketConnection(conn, session),
	}
}
