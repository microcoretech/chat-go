package http

import (
	"errors"
	"strings"

	"golang.org/x/exp/slices"

	"mbobrovskyi/chat-go/internal/common/domain"
)

func SortFromDto(querySort string, sortFields []string) (*domain.Sort, error) {
	var sort *domain.Sort

	if len(querySort) > 0 {
		parts := strings.Split(querySort, ",")

		if len(parts) > 0 {
			sort = &domain.Sort{
				SortBy:  parts[0],
				SortDir: domain.Asc,
			}

			if !slices.Contains(sortFields, sort.SortBy) {
				return nil, errors.New("invalid sort field")
			}
		}

		if len(parts) > 1 {
			sort.SortDir = domain.SortDirection(strings.ToLower(parts[1]))
			if sort.SortDir != domain.Asc && sort.SortDir != domain.Desc {
				return nil, errors.New("invalid sort direction")
			}
		}
	}

	return sort, nil
}
