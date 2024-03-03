package http

import (
	"chat/internal/common/domain"
	"chat/internal/common/errors"
	"chat/internal/common/http"
	"chat/internal/user/common"
)

var userSortFields = []string{
	"id",
	"email",
	"username",
	"role",
	"firstName",
	"lastName",
	"createdAt",
	"updatedAt",
}

func UserFilterFromQuery(query UserQuery) (domain.UserFilter, error) {
	sort, err := http.SortFromDto(query.Sort, userSortFields)
	if err != nil {
		return domain.UserFilter{}, errors.NewBadRequestError(
			common.UserDomain, err, map[string]any{"sort": query.Sort})
	}

	return domain.UserFilter{
		IDs:       query.IDs,
		Emails:    query.Emails,
		UserNames: query.Usernames,
		Roles:     query.Roles,
		Search:    query.Search,
		Limit:     query.Limit,
		Offset:    query.Offset,
		Sort:      sort,
	}, nil
}
