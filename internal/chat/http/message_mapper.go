package http

import (
	"github.com/samber/lo"

	"mbobrovskyi/chat-go/internal/chat/domain"
	"mbobrovskyi/chat-go/internal/common/http"
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
