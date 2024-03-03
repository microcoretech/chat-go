package repository

import (
	"chat/internal/chat/domain"
)

func messageFromDto(dto messageDto) domain.Message {
	return domain.Message{
		ID:        dto.ID,
		Text:      dto.Text,
		Status:    dto.Status,
		CreatedBy: dto.CreatedBy,
		Creator:   dto.Creator,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}
}
