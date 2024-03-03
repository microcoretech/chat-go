package websocket

import (
	"chat/internal/chat/domain"
	"chat/internal/common/http"
	"github.com/samber/lo"
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
		ChatID:    message.ChatID,
		Creator:   creatorDto,
		CreatedBy: message.CreatedBy,
		CreatedAt: message.CreatedAt,
		UpdatedAt: message.UpdatedAt,
	}
}

func MessageFromCreateDto(message MessageDto) domain.Message {
	return domain.Message{
		Text: message.Text,
	}
}
