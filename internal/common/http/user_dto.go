// Copyright 2025 Mykhailo Bobrovskyi
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
	"time"

	"chat-go/internal/common/domain"
)

type UserDto struct {
	ID        uint64       `json:"id"`
	Email     string       `json:"email"`
	Username  string       `json:"username"`
	Role      uint8        `json:"role"`
	FirstName string       `json:"firstName"`
	LastName  string       `json:"lastName"`
	AboutMe   string       `json:"aboutMe"`
	Image     domain.Image `json:"image"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
}
