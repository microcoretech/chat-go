package domain

import (
	"chat/internal/common/domain"
	"time"
)

type Token struct {
	Token   string
	ExpIn   time.Duration
	ExpTime time.Time
	User    domain.User
}
