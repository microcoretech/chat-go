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
	"mbobrovskyi/chat-go/internal/common/domain"
)

func UserFromSignUpRequest(req SignUpRequest) domain.User {
	return domain.User{
		Email:     req.Email,
		Username:  req.Username,
		Role:      domain.UserRole,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}
}
