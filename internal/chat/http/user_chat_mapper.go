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

package http

import (
	"github.com/samber/lo"

	"chat-go/internal/chat/domain"
	"chat-go/internal/common/http"
)

func UserChatFromDto(userChatDto UserChatDto) domain.UserChat {
	return domain.UserChat{
		UserID: userChatDto.UserID,
		ChatID: userChatDto.ChatID,
	}
}

func UserChatToDto(userChat domain.UserChat) UserChatDto {
	var userDto *http.UserDto
	if userChat.User != nil {
		userDto = lo.ToPtr(http.UserToDto(*userChat.User))
	}

	return UserChatDto{
		UserID: userChat.UserID,
		ChatID: userChat.ChatID,
		User:   userDto,
	}
}
