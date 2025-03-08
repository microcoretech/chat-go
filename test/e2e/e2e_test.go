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

package e2e

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	chatdomain "chat-go/internal/chat/domain"
	chathttp "chat-go/internal/chat/http"
	"chat-go/internal/common/common"
	commonhttp "chat-go/internal/common/http"
	"chat-go/internal/infrastructure/api"
	"chat-go/test/util"
)

var _ = ginkgo.Describe("Chat", func() {
	var client util.HTTPClient

	ginkgo.BeforeEach(func() {
		client = &http.Client{
			Timeout: util.Timeout,
		}

		util.RemoveAllChats(client, chatURL, util.AdminToken)
	})

	ginkgo.Context("root endpoint", func() {
		ginkgo.It("should return valid response", func() {
			req, err := http.NewRequest(http.MethodGet, chatURL, nil)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())

			resp, err := client.Do(req)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(resp.StatusCode).To(gomega.Equal(fiber.StatusOK))

			body, err := io.ReadAll(resp.Body)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			gotRootResponse := &api.RootResponse{}
			gomega.Expect(json.Unmarshal(body, gotRootResponse)).To(gomega.Succeed())

			projectDir, err := util.GetProjectDir()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			fileVersion, err := os.ReadFile(filepath.Join(projectDir, "VERSION"))
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			wantRootResponse := &api.RootResponse{
				Service: common.ServiceName,
				Version: string(fileVersion),
			}
			gomega.Expect(gotRootResponse).To(gomega.BeComparableTo(wantRootResponse))
		})
	})

	ginkgo.Context("health endpoint", func() {
		ginkgo.It("should return valid response", func() {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/healthz", chatURL), nil)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())

			resp, err := client.Do(req)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(resp.StatusCode).To(gomega.Equal(fiber.StatusOK))

			body, err := io.ReadAll(resp.Body)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(string(body)).To(gomega.Equal("OK"))
		})
	})

	ginkgo.Context("create chat endpoint", func() {
		ginkgo.It("should return valid response for group chat", func() {
			createChatRequest := &chathttp.CreateChatDto{
				Name: "Test Group Chat",
				Type: uint8(chatdomain.GroupChatType),
			}

			gomega.Expect(util.CreateChat(client, chatURL, util.AdminToken, createChatRequest)).To(gomega.BeComparableTo(
				&chathttp.ChatDto{
					Name:      "Test Group Chat",
					Type:      uint8(chatdomain.GroupChatType),
					CreatedBy: util.AdminID,
					Creator: &commonhttp.UserDto{
						ID:       util.AdminID,
						Email:    util.AdminEmail,
						Username: util.AdminUsername,
					},
					UserChats: []chathttp.UserChatDto{
						{
							UserID: util.AdminID,
							User: &commonhttp.UserDto{
								ID:       util.AdminID,
								Email:    util.AdminEmail,
								Username: util.AdminUsername,
							},
						},
					},
				},
				cmpopts.IgnoreFields(chathttp.ChatDto{}, "ID", "CreatedAt", "UpdatedAt"),
				cmpopts.IgnoreFields(chathttp.UserChatDto{}, "ChatID"),
			))
		})
	})

	ginkgo.Context("update chat endpoint", func() {
		ginkgo.It("should return valid response for updating group chat", func() {
			createChatRequest := &chathttp.CreateChatDto{
				Name: "Test Group Chat",
				Type: uint8(chatdomain.GroupChatType),
			}
			createdChat := util.CreateChat(client, chatURL, util.AdminToken, createChatRequest)

			updateChatRequest := &chathttp.UpdateChatDto{
				Name: "Updated Group Chat",
			}

			expectedChatResponse := chathttp.ChatDto{
				ID:        createdChat.ID,
				Name:      "Updated Group Chat",
				Type:      uint8(chatdomain.GroupChatType),
				CreatedBy: util.AdminID,
				UserChats: []chathttp.UserChatDto{
					{
						UserID: util.AdminID,
						ChatID: createdChat.ID,
					},
				},
			}

			updatedChat := util.UpdateChat(client, chatURL, util.AdminToken, createdChat.ID, updateChatRequest, http.StatusOK)
			gomega.Expect(updatedChat).To(gomega.BeComparableTo(
				&expectedChatResponse,
				cmpopts.IgnoreFields(chathttp.ChatDto{}, "CreatedAt", "UpdatedAt"),
			))
			gomega.Expect(updatedChat.UpdatedAt).ToNot(gomega.BeZero())
			gomega.Expect(updatedChat.UpdatedAt.After(updatedChat.CreatedAt) || updatedChat.UpdatedAt.Equal(updatedChat.CreatedAt)).To(gomega.BeTrue(), "UpdatedAt should be after or equal to CreatedAt")
		})
	})
})
