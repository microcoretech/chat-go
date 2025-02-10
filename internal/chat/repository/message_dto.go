package repository

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"mbobrovskyi/chat-go/internal/chat/domain"
)

type messageDto domain.Message

func (uc messageDto) Value() (driver.Value, error) {
	return json.Marshal(uc)
}

func (uc *messageDto) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &uc)
}
