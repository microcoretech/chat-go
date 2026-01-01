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
	"time"

	"chat-go/internal/common/http"
)

type MessageDto struct {
	UUID      string        `json:"uuid"`
	ID        uint64        `json:"id"`
	Text      string        `json:"text"`
	Status    uint8         `json:"status"`
	ChatID    uint64        `json:"chatId"`
	Creator   *http.UserDto `json:"creator"`
	CreatedBy uint64        `json:"createdBy"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
}
