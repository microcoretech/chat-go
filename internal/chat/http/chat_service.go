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

package http

import (
	"context"

	"chat-go/internal/chat/domain"
)

type ChatService interface {
	GetChat(ctx context.Context, id uint64) (*domain.Chat, error)
	GetChats(ctx context.Context, filter *domain.ChatFilter) ([]domain.Chat, uint64, error)
	CreateChat(ctx context.Context, chat domain.Chat) (*domain.Chat, error)
	UpdateChat(ctx context.Context, chat domain.Chat) (*domain.Chat, error)
	DeleteChat(ctx context.Context, id uint64) error
}
