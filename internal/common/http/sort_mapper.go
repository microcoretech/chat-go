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
	"errors"
	"strings"

	"golang.org/x/exp/slices"

	"chat-go/internal/common/domain"
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
