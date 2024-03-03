package repository

import (
	"chat/internal/chat/domain"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type userChatsDto []domain.UserChat

func (uc userChatsDto) Value() (driver.Value, error) {
	return json.Marshal(uc)
}

func (uc *userChatsDto) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &uc)
}
