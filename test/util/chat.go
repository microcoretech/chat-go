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

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/onsi/gomega"

	chathttp "chat-go/internal/chat/http"
	commonhttp "chat-go/internal/common/http"
)

func GetChats(client HTTPClient, baseURL string, token string) commonhttp.Page[chathttp.ChatDto] {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/chats", baseURL), nil)
	gomega.ExpectWithOffset(1, err).NotTo(gomega.HaveOccurred())

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	gomega.ExpectWithOffset(1, err).ToNot(gomega.HaveOccurred())
	defer resp.Body.Close()

	gomega.ExpectWithOffset(1, resp.StatusCode).To(gomega.Equal(fiber.StatusOK))

	responseBody, err := io.ReadAll(resp.Body)
	gomega.ExpectWithOffset(1, err).ToNot(gomega.HaveOccurred())

	var chats commonhttp.Page[chathttp.ChatDto]
	gomega.ExpectWithOffset(1, json.Unmarshal(responseBody, &chats)).To(gomega.Succeed())

	return chats
}

func CreateChat(client HTTPClient, baseURL string, token string, createChatRequest *chathttp.CreateChatDto) *chathttp.ChatDto {
	requestBody, err := json.Marshal(createChatRequest)
	gomega.ExpectWithOffset(1, err).NotTo(gomega.HaveOccurred())

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/chats", baseURL), bytes.NewBuffer(requestBody))
	gomega.ExpectWithOffset(1, err).ToNot(gomega.HaveOccurred())

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	gomega.ExpectWithOffset(1, err).ToNot(gomega.HaveOccurred())
	defer resp.Body.Close()

	gomega.ExpectWithOffset(1, resp.StatusCode).To(gomega.Equal(http.StatusOK))

	responseBody, err := io.ReadAll(resp.Body)
	gomega.ExpectWithOffset(1, err).ToNot(gomega.HaveOccurred())

	var createChatResponse chathttp.ChatDto
	gomega.ExpectWithOffset(1, json.Unmarshal(responseBody, &createChatResponse)).To(gomega.Succeed())

	return &createChatResponse
}

func DeleteChat(client HTTPClient, baseURL string, token string, id uint64) {
	deleteChat(client, baseURL, token, id, http.StatusOK)
}

func DeleteChatWithStatus(client HTTPClient, baseURL string, token string, id uint64, status int) {
	deleteChat(client, baseURL, token, id, status)
}

func deleteChat(client HTTPClient, baseURL string, token string, id uint64, status int) {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/chats/%d", baseURL, id), nil)
	gomega.ExpectWithOffset(2, err).ToNot(gomega.HaveOccurred())

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	gomega.ExpectWithOffset(2, err).ToNot(gomega.HaveOccurred())
	defer resp.Body.Close()

	gomega.ExpectWithOffset(2, resp.StatusCode).To(gomega.Equal(status))
}
