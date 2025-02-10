package http

import (
	"mbobrovskyi/chat-go/internal/chat/common"
	"mbobrovskyi/chat-go/internal/chat/domain"
	"mbobrovskyi/chat-go/internal/common/errors"
	"mbobrovskyi/chat-go/internal/common/http"
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
