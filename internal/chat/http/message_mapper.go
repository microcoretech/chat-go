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
	"github.com/samber/lo"

	"chat-go/internal/chat/domain"
	"chat-go/internal/common/http"
)

func MessageToDto(message domain.Message) MessageDto {
	var creatorDto *http.UserDto
	if message.Creator != nil {
		creatorDto = lo.ToPtr(http.UserToDto(*message.Creator))
	}

	return MessageDto{
		ID:        message.ID,
		Text:      message.Text,
		Status:    message.Status.ToUint8(),
		Creator:   creatorDto,
		CreatedBy: message.CreatedBy,
		CreatedAt: message.CreatedAt,
		UpdatedAt: message.UpdatedAt,
	}
}
