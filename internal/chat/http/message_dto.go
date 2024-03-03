package http

import (
	"chat/internal/common/http"
	"time"
)

type MessageDto struct {
	ID        uint64        `json:"id"`
	Text      string        `json:"text"`
	Status    uint8         `json:"status"`
	CreatedBy uint64        `json:"createdBy"`
	Creator   *http.UserDto `json:"creator"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
}
