package http

import (
	"chat/internal/chat/common"
	"chat/internal/chat/domain"
	"chat/internal/common/errors"
	"chat/internal/common/http"
)

var chatSortFields = []string{
	"id",
	"name",
	"type",
	"createdBy",
	"createdAt",
	"updatedAt",
}

func ChatFilterFromQuery(query ChatQuery) (domain.ChatFilter, error) {
	sort, err := http.SortFromDto(query.Sort, chatSortFields)
	if err != nil {
		return domain.ChatFilter{}, errors.NewBadRequestError(
			common.ChatDomain, err, map[string]any{"sort": query.Sort})
	}

	return domain.ChatFilter{
		IDs:          query.IDs,
		Types:        query.Types,
		CreatedByIDs: query.CreatedByIDs,
		Search:       query.Search,
		Limit:        query.Limit,
		Offset:       query.Offset,
		Sort:         sort,
	}, nil
}
