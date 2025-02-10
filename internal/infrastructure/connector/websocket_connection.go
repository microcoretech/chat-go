package connector

import (
	"encoding/json"

	"github.com/fasthttp/websocket"
	"github.com/google/uuid"

	"mbobrovskyi/chat-go/internal/common/domain"
)

type WebsocketConnection struct {
	connectionID string

	conn      *websocket.Conn
	connector Connector
	session   *domain.Session

	messageChan chan []byte
	closeChan   chan struct{}

	isConnected bool
	isClosed    bool
}

func (c *WebsocketConnection) IsClosed() bool {
	return c.isClosed
}

func (c *WebsocketConnection) GetConnectionID() string {
	return c.connectionID
}

func (c *WebsocketConnection) GetConnector() Connector {
	return c.connector
}

func (c *WebsocketConnection) SetConnector(connector Connector) {
	c.connector = connector
}

func (c *WebsocketConnection) GetMessageChan() chan []byte {
	return c.messageChan
}

func (c *WebsocketConnection) GetCloseChan() chan struct{} {
	return c.closeChan
}

func (c *WebsocketConnection) SendEvent(eventType uint64, data any) error {
	var err error

	event := Event{
		Type: eventType,
	}

	event.Data, err = json.Marshal(data)
	if err != nil {
		return err
	}

	err = c.conn.WriteJSON(event)
	if err != nil {
		return err
	}

	return nil
}

func (c *WebsocketConnection) Connect() {
	go c.connect()
}

func (c *WebsocketConnection) connect() {
	if c.isConnected || c.isClosed {
		return
	}

	c.isConnected = true

	defer func() {
		c.isConnected = false
		c.isClosed = true
	}()

	for {
		select {
		case <-c.closeChan:
			return
		default:
			_, msgData, err := c.conn.ReadMessage()
			if err != nil {
				c.isConnected = false
				return
			}
			c.messageChan <- msgData
		}
	}
}

func (c *WebsocketConnection) Close() {
	if !c.isConnected {
		return
	}
	close(c.closeChan)
	c.conn.Close()
}

func (c *WebsocketConnection) GetSession() *domain.Session {
	return c.session
}

func NewWebSocketConnection(conn *websocket.Conn, session *domain.Session) *WebsocketConnection {
	return &WebsocketConnection{
		connectionID: uuid.NewString(),
		conn:         conn,
		session:      session,
		messageChan:  make(chan []byte),
		closeChan:    make(chan struct{}),
	}
}
