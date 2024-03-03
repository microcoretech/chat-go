package domain

type UserFilter struct {
	IDs       []uint64
	Emails    []string
	UserNames []string
	Roles     []uint8

	Search string

	Limit  *uint64
	Offset *uint64

	Sort *Sort
}
