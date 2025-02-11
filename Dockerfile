# Copyright 2025 Mykhailo Bobrovskyi
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM golang:1.23-alpine AS builder

RUN apk --no-cache add ca-certificates git

WORKDIR /go/src/chat-go

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./bin/chat-go ./cmd/chat

FROM scratch
WORKDIR /
COPY --from=builder /go/src/chat-go/bin/chat-go /bin/chat-go
COPY /VERSION .
ENTRYPOINT ["/bin/chat-go"]