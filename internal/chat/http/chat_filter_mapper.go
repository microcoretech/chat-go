// Copyright 2025 MicroCore Tech
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package http

import (
	"chat-go/internal/chat/constants"
	"chat-go/internal/chat/domain"
	"chat-go/internal/common/errors"
	"chat-go/internal/common/http"
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
			constants.ChatDomain, err, map[string]any{"sort": query.Sort})
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
