package http

import (
	"chat/internal/chat/common"
	"chat/internal/chat/domain"
	"chat/internal/common/errors"
	"chat/internal/common/http"
)

var messageSortFields = []string{
	"id",
	"status",
	"chatId",
	"createdBy",
	"createdAt",
	"updatedAt",
}

func MessageFilterFromQuery(query MessageQuery) (domain.MessageFilter, error) {
	sort, err := http.SortFromDto(query.Sort, messageSortFields)
	if err != nil {
		return domain.MessageFilter{}, errors.NewBadRequestError(
			common.ChatDomain, err, map[string]any{"sort": query.Sort})
	}

	return domain.MessageFilter{
		IDs:          query.IDs,
		ChatIDs:      query.ChatIDs,
		Statuses:     query.Statuses,
		CreatedByIDs: query.CreatedByIDs,
		Search:       query.Search,
		Limit:        query.Limit,
		Offset:       query.Offset,
		Sort:         sort,
	}, nil
}
