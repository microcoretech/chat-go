package domain

import (
	"github.com/google/uuid"
)

type Session struct {
	ID   uuid.UUID `json:"id"`
	User User      `json:"user"`
}
