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

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	chat "chat-go/internal/chat/http"
	"chat-go/internal/common/domain"
	user "chat-go/internal/common/http"

	"chat-go/internal/common/common"
	"chat-go/internal/infrastructure/api"
	"chat-go/test/util"
)

var _ = ginkgo.Describe("Chat", func() {
	var client *http.Client
	var adminToken string
	var adminUserId uint64

	ginkgo.BeforeEach(func() {
		client = &http.Client{
			Timeout: 5 * time.Second,
		}

		ginkgo.DeferCleanup(func() {
			client.CloseIdleConnections()
		})

		signInRequest := map[string]any{
			"password": "mohnIeih4Ju9zHYE1VPWL0mHyzBjyFPl",
			"username": "admin",
		}

		reqBody, err := json.Marshal(signInRequest)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		signInReq, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/auth/sign-in", chatURL), bytes.NewBuffer(reqBody))
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		signInReq.Header.Set("Content-Type", "application/json")

		signInResp, err := client.Do(signInReq)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		body, err := io.ReadAll(signInResp.Body)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		defer signInResp.Body.Close()

		gomega.Expect(signInResp.StatusCode).To(gomega.Equal(http.StatusOK))

		var tokenResponse map[string]interface{}
		gomega.Expect(json.Unmarshal(body, &tokenResponse)).To(gomega.Succeed())
		adminToken = tokenResponse["token"].(string)

		user, ok := tokenResponse["user"].(map[string]interface{})
		if !ok {
			ginkgo.Fail("Sign-in response missing 'user' field")
		}
		adminUserId = uint64(user["id"].(float64))
	})

	ginkgo.Context("root endpoint", func() {
		ginkgo.It("should return valid response", func() {
			resp, err := client.Get(chatURL)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer resp.Body.Close()
			gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK))

			body, err := io.ReadAll(resp.Body)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			gotRootResponse := &api.RootResponse{}
			gomega.Expect(json.Unmarshal(body, gotRootResponse)).To(gomega.Succeed())

			fileVersion, err := os.ReadFile(filepath.Join(util.ProjectDir, "VERSION"))
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
			defer resp.Body.Close()
			gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK))

			body, err := io.ReadAll(resp.Body)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(string(body)).To(gomega.Equal("OK"))
		})
	})

	//  TODO ADD CLEANUP TEST

	ginkgo.It("should return valid response for create group chat", func() {
		createChatRequest := chat.CreateChatDto{
			Name: "Test Group Chat",
			Type: 2,
		}
		chatReqBody, err := json.Marshal(createChatRequest)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/chats", chatURL), bytes.NewBuffer(chatReqBody))
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		body, err := io.ReadAll(resp.Body)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		defer resp.Body.Close()

		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK))

		var chatResponse chat.ChatDto
		gomega.Expect(json.Unmarshal(body, &chatResponse)).To(gomega.Succeed())

		expectedChatResponse := chat.ChatDto{
			ID:          chatResponse.ID,
			Name:        "Test Group Chat",
			Type:        2,
			Image:       domain.Image{},
			LastMessage: nil,
			CreatedBy:   adminUserId,
			Creator: &user.UserDto{
				ID:        adminUserId,
				Email:     "admin@gmail.com",
				Username:  "admin",
				Role:      2,
				FirstName: "Admin",
				LastName:  "Admin",
				AboutMe:   "",
				Image:     domain.Image{},
				CreatedAt: chatResponse.Creator.CreatedAt,
				UpdatedAt: chatResponse.Creator.UpdatedAt,
			},
			UserChats: []chat.UserChatDto{
				{
					UserID: adminUserId,
					ChatID: chatResponse.ID,
					User: &user.UserDto{
						ID:        adminUserId,
						Email:     "admin@gmail.com",
						Username:  "admin",
						Role:      2,
						FirstName: "Admin",
						LastName:  "Admin",
						AboutMe:   "",
						Image:     domain.Image{},
						CreatedAt: chatResponse.UserChats[0].User.CreatedAt,
						UpdatedAt: chatResponse.UserChats[0].User.UpdatedAt,
					},
				},
			},
			CreatedAt: chatResponse.CreatedAt,
			UpdatedAt: chatResponse.UpdatedAt,
		}

		gomega.Expect(chatResponse).To(gomega.BeComparableTo(expectedChatResponse))
		gomega.Expect(chatResponse.CreatedAt).ToNot(gomega.BeZero())
		gomega.Expect(chatResponse.UpdatedAt).ToNot(gomega.BeZero())
		gomega.Expect(chatResponse.Creator.CreatedAt).ToNot(gomega.BeZero())
		gomega.Expect(chatResponse.Creator.UpdatedAt).ToNot(gomega.BeZero())
		gomega.Expect(chatResponse.UserChats[0].User.CreatedAt).ToNot(gomega.BeZero())
		gomega.Expect(chatResponse.UserChats[0].User.UpdatedAt).ToNot(gomega.BeZero())
	})
})
