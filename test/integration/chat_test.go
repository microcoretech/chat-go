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

package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/gofiber/fiber/v2"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	chatdomain "chat-go/internal/chat/domain"
	chathttp "chat-go/internal/chat/http"
	commonhttp "chat-go/internal/common/http"
	"chat-go/test/util"
)

var _ = ginkgo.Describe("Chat", func() {
	var (
		chat *chathttp.ChatDto
	)

	ginkgo.BeforeEach(func() {
		chat = createChat(&chathttp.CreateChatDto{
			Name: "Chat1",
			Type: uint8(chatdomain.GroupChatType),
		})
	})

	ginkgo.AfterEach(func() {
		// TODO: Cleanup created chat after implementation https://github.com/mbobrovskyi/chat-go/issues/28
	})

	ginkgo.Context("get chats endpoint", func() {
		ginkgo.It("should return valid response", func() {
			req := httptest.NewRequest(http.MethodGet, "/chats", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", util.AdminUsername))

			resp, err := env.App().Test(req)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(resp.StatusCode).To(gomega.Equal(fiber.StatusOK))

			responseBody, err := io.ReadAll(resp.Body)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())

			var chats commonhttp.Page[chathttp.ChatDto]
			gomega.Expect(json.Unmarshal(responseBody, &chats)).To(gomega.Succeed())

			gomega.Expect(chats.Items).To(gomega.HaveLen(1))
			gomega.Expect(chats.Count).To(gomega.Equal(uint64(1)))
			gomega.Expect(&chats.Items[0]).To(gomega.BeComparableTo(chat))
		})
	})
})

func createChat(createChatRequest *chathttp.CreateChatDto) *chathttp.ChatDto {
	requestBody, err := json.Marshal(createChatRequest)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	req := httptest.NewRequest(http.MethodPost, "/chats", bytes.NewBuffer(requestBody))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", util.AdminUsername))
	req.Header.Set("Content-Type", "application/json")

	resp, err := env.App().Test(req)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	gomega.Expect(resp.StatusCode).To(gomega.Equal(fiber.StatusOK))

	responseBody, err := io.ReadAll(resp.Body)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	var createChatResponse chathttp.ChatDto
	gomega.Expect(json.Unmarshal(responseBody, &createChatResponse)).To(gomega.Succeed())

	return &createChatResponse
}
