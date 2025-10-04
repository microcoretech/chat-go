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

package websocket

const (
	SubscribeChatsEventType       = 1
	UnsubscribeChatsEventType     = 2
	SetCurrentChatEventType       = 3
	UnsetCurrentChatEventType     = 4
	CreateMessageEventType        = 5
	EditMessageEventType          = 6
	DeleteMessageEventType        = 7
	UpdateMessagesStatusEventType = 8
)

type EditMessageEventData struct {
	MessageID string `json:"messageId"`
	Text      string `json:"text"`
}

type DeleteMessageEventData struct {
	MessageID string `json:"messageId"`
}
