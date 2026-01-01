// Copyright MicroCore Tech
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

package domain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/samber/lo"

	"chat-go/internal/common/domain"
	commonerrors "chat-go/internal/common/errors"
	commonhttp "chat-go/internal/common/http"
	"chat-go/internal/infrastructure/configs"
	"chat-go/internal/user/constants"
)

type UserServiceImpl struct {
	config *configs.Config
}

func (s *UserServiceImpl) GetCurrentUser(ctx context.Context) (*domain.User, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", s.config.GetCurrentUserEndpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", domain.TokenFromContext(ctx)))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return nil, commonerrors.NewUnauthorizedError()
	case http.StatusNotFound:
		return nil, commonerrors.NewUndefinedError(
			errors.New("invalid get current user endpoint"),
			"endpoint", s.config.GetCurrentUserEndpoint,
		)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	userDto := commonhttp.UserDto{}
	err = json.Unmarshal(body, &userDto)
	if err != nil {
		return nil, err
	}

	if userDto.ID == 0 {
		return nil, commonerrors.NewUnauthorizedError()
	}

	return lo.ToPtr(commonhttp.UserFromDto(userDto)), nil
}

func (s *UserServiceImpl) GetUsers(ctx context.Context, filter *domain.UserFilter) ([]domain.User, uint64, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", s.config.GetUsersEndpoint, nil)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", domain.TokenFromContext(ctx)))

	q := req.URL.Query()

	for _, id := range filter.IDs {
		q.Add("ids", strconv.FormatUint(id, 10))
	}

	for _, email := range filter.Emails {
		q.Add("emails", email)
	}

	for _, username := range filter.Usernames {
		q.Add("usernames", username)
	}

	if filter.Search != "" {
		q.Add("search", filter.Search)
	}

	if filter.Limit != nil {
		q.Add("limit", fmt.Sprint(filter.Limit))
	}

	if filter.Offset != nil {
		q.Add("offset", fmt.Sprint(filter.Offset))
	}

	if filter.Sort != nil {
		q.Add("sort", fmt.Sprintf("%s,%s", filter.Sort.SortBy, filter.Sort.SortDir.String()))
	}

	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	switch resp.StatusCode {
	case http.StatusBadRequest:
		return nil, 0, commonerrors.NewBadRequestError(constants.UserDomain, nil, map[string]any{
			"body": string(body),
		})
	case http.StatusUnauthorized:
		return nil, 0, commonerrors.NewUnauthorizedError()
	}

	pageDto := &commonhttp.Page[commonhttp.UserDto]{}
	err = json.Unmarshal(body, pageDto)
	if err != nil {
		return nil, 0, err
	}

	return lo.Map(pageDto.Items, func(userDto commonhttp.UserDto, _ int) domain.User {
		return commonhttp.UserFromDto(userDto)
	}), pageDto.Count, nil
}

func NewUserServiceImpl(config *configs.Config) *UserServiceImpl {
	return &UserServiceImpl{
		config: config,
	}
}
