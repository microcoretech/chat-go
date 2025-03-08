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
	"github.com/fasthttp/websocket"
	"golang.org/x/exp/slices"

	"chat-go/internal/infrastructure/connector"
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

func NewConnection(conn *websocket.Conn) connector.Connection {
	return &connectionImpl{
		Connection: connector.NewWebSocketConnection(conn),
	}
}
