package http

type MessageQuery struct {
	IDs          []uint64 `query:"id" validate:"omitempty,gte=0"`
	ChatIDs      []uint64 `query:"chatId" validate:"omitempty,gte=0"`
	Statuses     []uint8  `query:"statuses" validate:"omitempty,oneof=1 2 3"`
	CreatedByIDs []uint64 `query:"id" validate:"omitempty,gte=0"`

	Search string `query:"search"`

	Limit  *uint64 `query:"limit"`
	Offset *uint64 `query:"offset"`

	Sort string `query:"sort"`
}
