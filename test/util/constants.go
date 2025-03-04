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

package util

import "time"

const (
	ChatImage       = "chat-go"
	PostgresImage   = "postgres:17-alpine"
	RedisImage      = "redis:7-alpine"
	MockServerImage = "mockserver/mockserver:5.15.0"
)

const (
	Timeout = time.Second * 10
)

const (
	AdminID       = 1
	AdminEmail    = "admin@gmail.com"
	AdminUsername = "admin"
	UserID        = 2
	UserEmail     = "user@gmail.com"
	UserUsername  = "user"
)
