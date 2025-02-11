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
	"golang.org/x/exp/slices"

	"chat-go/internal/chat/errors"
)

type ChatType uint8

const (
	DirectChatType ChatType = 1
	GroupChatType  ChatType = 2
)

func (ct ChatType) Uint8() uint8 {
	return uint8(ct)
}

func (ct ChatType) Types() []ChatType {
	return []ChatType{DirectChatType, GroupChatType}
}

func (ct ChatType) IsValid() bool {
	return slices.Contains(ct.Types(), ct)
}

func NewChatType(chatType uint8) (ChatType, error) {
	if !ChatType(chatType).IsValid() {
		return 0, errors.NewInvalidChatTypeError()
	}

	return ChatType(chatType), nil
}
