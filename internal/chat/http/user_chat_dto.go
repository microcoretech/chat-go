package http

import (
	"chat/internal/common/http"
)

type UserChatDto struct {
	UserID uint64 `json:"userId"`
	ChatID uint64 `json:"chatId"`

	User *http.UserDto `json:"user"`
}
