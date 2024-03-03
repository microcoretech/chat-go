package domain

import (
	"github.com/google/uuid"
)

func UserToSession(user User) Session {
	return Session{
		ID:   uuid.New(),
		User: user,
	}
}
