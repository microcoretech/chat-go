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
	"context"

	"chat-go/internal/common/domain"
	userdomain "chat-go/internal/user/domain"
)

type AuthService interface {
	SignIn(ctx context.Context, username, password string) (*userdomain.Token, error)
	SignUp(ctx context.Context, newUser domain.User, password string) (*userdomain.Token, error)
	SignOut(ctx context.Context, userID uint64, token string) error
	GetSession(ctx context.Context, token string) (*domain.Session, error)
}
