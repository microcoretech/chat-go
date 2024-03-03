package domain

import (
	"time"
)

const (
	UserRole  uint8 = 1
	AdminRole uint8 = 2
)

type User struct {
	ID        uint64    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Role      uint8     `json:"role"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	AboutMe   string    `json:"aboutMe"`
	Image     Image     `json:"image"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
