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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

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
	var client *http.Client

	ginkgo.BeforeEach(func() {
		client = &http.Client{
			Timeout: 5 * time.Second,
		}

		ginkgo.DeferCleanup(func() {
			client.CloseIdleConnections()
		})
	})

	ginkgo.Context("root endpoint", func() {
		ginkgo.It("should return valid response", func() {
			resp, err := client.Get(chatURL)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK))

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
			resp, err := client.Get(fmt.Sprintf("%s/healthz", chatURL))
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK))

			body, err := io.ReadAll(resp.Body)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(string(body)).To(gomega.Equal("OK"))
		})
	})

	ginkgo.Context("chat create endpoint", func() {
		// TODO: Cleanup created chat after implementation https://github.com/mbobrovskyi/chat-go/issues/28
		ginkgo.It("should return valid response for create group chat", func() {
			createChatRequest := chathttp.CreateChatDto{
				Name: "Test Group Chat",
				Type: uint8(chatdomain.GroupChatType),
			}
			chatReqBody, err := json.Marshal(createChatRequest)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/chats", chatURL), bytes.NewBuffer(chatReqBody))
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", util.AdminUsername))
			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			body, err := io.ReadAll(resp.Body)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK))

			var chatResponse chathttp.ChatDto
			gomega.Expect(json.Unmarshal(body, &chatResponse)).To(gomega.Succeed())

			expectedChatResponse := chathttp.ChatDto{
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
			}

			gomega.Expect(chatResponse).To(gomega.BeComparableTo(
				expectedChatResponse,
				cmpopts.IgnoreFields(chathttp.ChatDto{}, "ID", "CreatedAt", "UpdatedAt"),
				cmpopts.IgnoreFields(chathttp.UserChatDto{}, "ChatID"),
			))
		})
	})
})
