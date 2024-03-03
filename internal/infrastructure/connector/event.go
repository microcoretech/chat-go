package connector

import (
	"encoding/json"
)

type Event struct {
	Type uint64          `json:"type"`
	Data json.RawMessage `json:"data"`
}
