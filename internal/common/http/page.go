package http

type Page[T any] struct {
	Items []T    `json:"items"`
	Count uint64 `json:"count"`
}

func NewPage[T any](items []T, count uint64) Page[T] {
	return Page[T]{Items: items, Count: count}
}
