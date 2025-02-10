package http

import (
	"time"

	"mbobrovskyi/chat-go/internal/common/http"
)

type TokenDto struct {
	Token   string       `json:"token"`
	ExpIn   int64        `json:"expIn"`
	ExpTime time.Time    `json:"expTime"`
	User    http.UserDto `json:"user"`
}
