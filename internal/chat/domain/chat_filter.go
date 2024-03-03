package domain

import "chat/internal/common/domain"

type ChatFilter struct {
	IDs          []uint64
	Types        []uint8
	CreatedByIDs []uint64

	Search string

	Limit  *uint64
	Offset *uint64

	Sort *domain.Sort
}
