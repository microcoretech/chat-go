package websocket

import (
	"time"

	"mbobrovskyi/chat-go/internal/common/http"
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
