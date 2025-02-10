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

package domain

import (
	"context"

	"mbobrovskyi/chat-go/internal/common/repository"
)

type ChatRepo interface {
	GetChat(ctx context.Context, id uint64) (*Chat, error)
	GetChats(ctx context.Context, filter *ChatFilter) ([]Chat, error)
	GetChatsCount(ctx context.Context, filter *ChatFilter) (uint64, error)
	CreateChat(ctx context.Context, chat Chat, tx repository.Tx) (*Chat, error)
	UpdateChat(ctx context.Context, chat Chat) (*Chat, error)
	DeleteChat(ctx context.Context, id uint64) (bool, error)
}
