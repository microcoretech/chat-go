package domain

import "mbobrovskyi/chat-go/internal/common/domain"

type MessageFilter struct {
	IDs          []uint64
	ChatIDs      []uint64
	Statuses     []uint8
	CreatedByIDs []uint64

	Search string

	Limit  *uint64
	Offset *uint64

	Sort *domain.Sort
}
