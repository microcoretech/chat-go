// Copyright 2025 MicroCore Tech
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
	"context"
	"encoding/json"
	"errors"
	"runtime/debug"
	"sync"
	"time"

	"chat-go/internal/infrastructure/logger"
)

var ErrConnectorAlreadyStarted = errors.New("connector already started")

type Connector interface {
	Start(ctx context.Context) error
	AddConnection(conn Connection)
	GetConnections() []Connection
}

type ConnectorImpl struct {
	mtx          sync.RWMutex
	log          logger.Logger
	connections  []Connection
	isStarted    bool
	eventHandler EventHandler
}

func (c *ConnectorImpl) Start(ctx context.Context) error {
	if c.isStarted {
		return ErrConnectorAlreadyStarted
	}

	c.isStarted = true
	defer func() {
		c.isStarted = false
	}()

	for {
		select {
		case <-ctx.Done():
			c.closeAll()
			return nil
		case <-time.After(time.Minute):
			c.log.Debug("Clean closed connections")
			c.clean()
		}
	}
}

func (c *ConnectorImpl) closeAll() {
	for _, conn := range c.connections {
		conn.Close()
	}
}

func (c *ConnectorImpl) clean() {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	connections := make([]Connection, 0)

	for _, conn := range c.connections {
		if !conn.IsClosed() {
			connections = append(connections, conn)
		}
	}

	c.connections = connections
}

func (c *ConnectorImpl) AddConnection(conn Connection) {
	c.log.Debugf("Connected id=%d email=%s username=%s", conn.GetUser().ID, conn.GetUser().Email, conn.GetUser().Username)
	conn.SetConnector(c)
	conn.Connect()
	c.addConnection(conn)
	go c.listen(conn)
}

func (c *ConnectorImpl) addConnection(conn Connection) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.connections = append(c.connections, conn)
}

func (c *ConnectorImpl) listen(conn Connection) {
	for {
		select {
		case <-conn.GetCloseChan():
			return
		case msg := <-conn.GetMessageChan():
			c.onEvent(conn, msg)
		}
	}
}

func (c *ConnectorImpl) onEvent(conn Connection, data []byte) {
	defer func() {
		if r := recover(); r != nil {
			c.log.Errorf("%s\n%s", r, string(debug.Stack()))
		}
	}()

	var rawEvent Event

	if err := json.Unmarshal(data, &rawEvent); err != nil {
		c.log.Debugf("error on parse raw event: %s", err.Error())
		return
	}

	c.log.Debugf("Got new event event_type=%d message=%s", rawEvent.Type, string(rawEvent.Data))

	if err := c.eventHandler.HandleEvent(conn, rawEvent); err != nil {
		c.log.Error(err)
	}
}

func (c *ConnectorImpl) GetConnections() []Connection {
	return c.connections
}

func NewConnector(
	log logger.Logger,
	eventHandler EventHandler,
) *ConnectorImpl {
	return &ConnectorImpl{
		log:          log,
		eventHandler: eventHandler,
	}
}
