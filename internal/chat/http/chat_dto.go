package http

import (
	"time"

	"mbobrovskyi/chat-go/internal/common/domain"
	"mbobrovskyi/chat-go/internal/common/http"
)

type CreateChatDto struct {
	Name      string        `json:"name" validate:"lte=255"`
	Type      uint8         `json:"type" validate:"required,oneof=1 2"`
	Image     domain.Image  `json:"image"`
	UserChats []UserChatDto `json:"users" validate:"dive,gte=0"`
}

type UpdateChatDto struct {
	ID    uint64       `json:"id" validate:"omitempty,gte=0"`
	Name  string       `json:"name" validate:"lte=255"`
	Type  uint8        `json:"type" validate:"required,oneof=1 2"`
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
