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

package http

import (
	"time"

	"chat-go/internal/common/domain"
	"chat-go/internal/common/http"
)

type CreateChatDto struct {
	Name      string        `json:"name" validate:"lte=255"`
	Type      uint8         `json:"type" validate:"required,oneof=1 2"`
	Image     domain.Image  `json:"image"`
	UserChats []UserChatDto `json:"users" validate:"dive,gte=0"`
}

type UpdateChatDto struct {
	Name  string       `json:"name" validate:"lte=255"`
	Image domain.Image `json:"image"`
}

type ChatDto struct {
	ID          uint64        `json:"id"`
	Name        string        `json:"name"`
	Type        uint8         `json:"type"`
	Image       domain.Image  `json:"image"`
	LastMessage *MessageDto   `json:"lastMessage"`
	CreatedBy   uint64        `json:"createdBy"`
	Creator     *http.UserDto `json:"creator"`
	UserChats   []UserChatDto `json:"userChats"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
}
