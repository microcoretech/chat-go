package http

import (
	"mbobrovskyi/chat-go/internal/chat/common"
	"mbobrovskyi/chat-go/internal/chat/domain"
	"mbobrovskyi/chat-go/internal/common/errors"
	"mbobrovskyi/chat-go/internal/common/http"
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
