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
	"net/http"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	chatdomain "chat-go/internal/chat/domain"
	chathttp "chat-go/internal/chat/http"
	"chat-go/test/util"
)

var _ = ginkgo.Describe("Chat", ginkgo.Ordered, ginkgo.ContinueOnFailure, func() {
	var (
		httpClient util.HTTPClient
		chat       *chathttp.ChatDto
	)

	ginkgo.BeforeAll(func() {
		httpClient = NewTestHTTPClient(env).WithTimeout(util.Timeout)

		chat = util.CreateChat(httpClient, "", util.AdminToken, &chathttp.CreateChatDto{
			Name: "Chat",
			Type: uint8(chatdomain.GroupChatType),
		})
	})

	ginkgo.AfterAll(func() {
		util.DeleteChat(httpClient, "", util.AdminToken, chat.ID)
	})

	ginkgo.Context("get chats endpoint", func() {
		ginkgo.It("should return valid response", func() {
			chats := util.GetChats(httpClient, "", util.AdminToken)
			gomega.Expect(chats.Items).To(gomega.HaveLen(1))
			gomega.Expect(chats.Count).To(gomega.Equal(uint64(1)))
			gomega.Expect(&chats.Items[0]).To(gomega.BeComparableTo(chat))
		})
	})

	ginkgo.Context("delete chat endpoint", func() {
		ginkgo.It("shouldn't delete not owned chat", func() {
			ginkgo.By("creating chat", func() {
				chat = util.CreateChat(httpClient, "", util.AdminToken, &chathttp.CreateChatDto{
					Name: "Chat",
					Type: uint8(chatdomain.GroupChatType),
				})
			})

			ginkgo.By("deleting chat", func() {
				util.DeleteChatWithStatus(httpClient, "", util.UserToken, chat.ID, http.StatusForbidden)
			})
		})
	})
})
