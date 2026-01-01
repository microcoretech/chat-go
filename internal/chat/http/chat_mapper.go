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

package http

import (
	"github.com/samber/lo"

	"chat-go/internal/chat/domain"
	"chat-go/internal/common/http"
)

func ChatFromCreateDto(dto CreateChatDto) (*domain.Chat, error) {
	chatType, err := domain.NewChatType(dto.Type)
	if err != nil {
		return nil, err
	}

	return &domain.Chat{
		Name:  dto.Name,
		Type:  chatType,
		Image: dto.Image,
		UserChats: lo.Map(dto.UserChats, func(userChat UserChatDto, _ int) domain.UserChat {
			return domain.UserChat{
				UserID: userChat.UserID,
				ChatID: userChat.ChatID,
			}
		}),
	}, nil
}

func ChatFromUpdateDto(dto UpdateChatDto) domain.Chat {
	return domain.Chat{
		Name:  dto.Name,
		Image: dto.Image,
	}
}

func ChatToDto(chat domain.Chat) ChatDto {
	var messageDto *MessageDto

	if chat.LastMessage != nil {
		messageDto = lo.ToPtr(MessageToDto(*chat.LastMessage))
	}

	var creator *http.UserDto
	if chat.Creator != nil {
		creator = lo.ToPtr(http.UserToDto(*chat.Creator))
	}

	return ChatDto{
		ID:          chat.ID,
		Name:        chat.Name,
		Type:        chat.Type.Uint8(),
		Image:       chat.Image,
		LastMessage: messageDto,
		CreatedBy:   chat.CreatedBy,
		Creator:     creator,
		UserChats: lo.Map(chat.UserChats, func(userChat domain.UserChat, _ int) UserChatDto {
			return UserChatToDto(userChat)
		}),
		CreatedAt: chat.CreatedAt,
		UpdatedAt: chat.UpdatedAt,
	}
}
