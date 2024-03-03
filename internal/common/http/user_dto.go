package http

import (
	"chat/internal/common/domain"
	"time"
)

type UserDto struct {
	ID        uint64       `json:"id"`
	Email     string       `json:"email"`
	Username  string       `json:"username"`
	Role      uint8        `json:"role"`
	FirstName string       `json:"firstName"`
	LastName  string       `json:"lastName"`
	AboutMe   string       `json:"aboutMe"`
	Image     domain.Image `json:"image"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
}
