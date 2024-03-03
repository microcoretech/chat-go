package domain

type SortDirection string

func (d SortDirection) String() string {
	return string(d)
}

const (
	Asc  SortDirection = "asc"
	Desc SortDirection = "desc"
)

type Sort struct {
	SortBy  string
	SortDir SortDirection
}
