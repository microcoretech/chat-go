package domain

import (
	"time"

	"mbobrovskyi/chat-go/internal/common/domain"
)

type Token struct {
	Token   string
	ExpIn   time.Duration
	ExpTime time.Time
	User    domain.User
}
