package http

import (
	"chat/internal/chat/domain"
	"chat/internal/common/http"
	"github.com/samber/lo"
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
