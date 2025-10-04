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

package errors

import (
	"encoding/json"
	"runtime/debug"
)

type BaseError interface {
	GetErrorData() *ErrorData
}

type ErrorData struct {
	Domain     string         `json:"domain"`
	ErrorType  string         `json:"type"`
	Data       map[string]any `json:"data"`
	DevDetails []string       `json:"devDetails"`
}

type ErrorDataShort struct {
	Domain    string         `json:"domain"`
	ErrorType string         `json:"type"`
	Data      map[string]any `json:"data"`
}

func TruncateErrorData(errorData *ErrorData) ErrorDataShort {
	delete(errorData.Data, "errorMessage")
	delete(errorData.Data, "stack")

	return ErrorDataShort{
		Domain:    errorData.Domain,
		ErrorType: errorData.ErrorType,
		Data:      errorData.Data,
	}
}

func (e *ErrorData) Error() string {
	data, _ := json.Marshal(e)
	return string(data)
}

func (e *ErrorData) GetErrorData() *ErrorData {
	return e
}

func NewErrorData(domain, errType string, err error, data map[string]any, devDetails ...string) *ErrorData {
	if data == nil {
		data = make(map[string]any)
	}

	if err != nil {
		stackStr := string(debug.Stack())
		data["stack"] = stackStr
		data["errorMessage"] = err.Error()
	}

	return &ErrorData{
		Domain:     domain,
		ErrorType:  errType,
		Data:       data,
		DevDetails: devDetails,
	}
}
