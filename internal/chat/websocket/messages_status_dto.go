package websocket

type MessagesStatusDto struct {
	Status     uint8    `json:"status" validate:"required,oneof=2 3"`
	MessageIDs []uint64 `json:"messageIds" validate:"required,gte=0"`
}
