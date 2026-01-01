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

package connector

import (
	"encoding/json"

	"github.com/fasthttp/websocket"
	"github.com/google/uuid"

	"chat-go/internal/common/domain"
)

type WebsocketConnection struct {
	connectionID string

	conn      *websocket.Conn
	connector Connector
	user      *domain.User

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

func (c *WebsocketConnection) GetUser() *domain.User {
	return c.user
}

func NewWebSocketConnection(conn *websocket.Conn, user *domain.User) *WebsocketConnection {
	return &WebsocketConnection{
		connectionID: uuid.NewString(),
		conn:         conn,
		user:         user,
		messageChan:  make(chan []byte),
		closeChan:    make(chan struct{}),
	}
}
