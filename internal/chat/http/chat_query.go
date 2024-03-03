package http

type ChatQuery struct {
	IDs          []uint64 `query:"id" validate:"omitempty,dive,gte=0"`
	Types        []uint8  `query:"types" validate:"omitempty,dive,oneof=1 2"`
	CreatedByIDs []uint64 `query:"createdByIds" validate:"omitempty,dive,gte=0"`

	Search string `query:"search"`

	Limit  *uint64 `query:"limit"`
	Offset *uint64 `query:"offset"`

	Sort string `query:"sort"`
}
