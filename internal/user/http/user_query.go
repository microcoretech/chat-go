package http

type UserQuery struct {
	IDs       []uint64 `query:"ids" validate:"omitempty,dive,gte=0"`
	Emails    []string `query:"emails" validate:"omitempty,dive,email"`
	Usernames []string `query:"usernames" validate:"omitempty,dive,gte=1,lte=255"`
	Roles     []uint8  `query:"roles" validate:"omitempty,dive,oneof=1 2"`

	Search string `query:"search"`

	Limit  *uint64 `query:"limit"`
	Offset *uint64 `query:"offset"`

	Sort string `query:"sort"`
}
