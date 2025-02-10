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

package errors

import (
	"mbobrovskyi/chat-go/internal/chat/common"
	"mbobrovskyi/chat-go/internal/common/errors"
)

const InvalidChatNameErrorType = "InvalidChatNameError"

type InvalidChatNameError struct {
	*errors.ErrorData
}

func NewInvalidChatNameError() *InvalidChatNameError {
	return &InvalidChatNameError{
		ErrorData: errors.NewErrorData(common.ChatDomain, InvalidChatNameErrorType, nil, nil),
	}
}
