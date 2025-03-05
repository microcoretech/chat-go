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

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	commonhttp "chat-go/internal/common/http"
	"chat-go/test/util"
)

var _ = ginkgo.Describe("User", func() {
	var client *http.Client

	ginkgo.BeforeEach(func() {
		client = &http.Client{
			Timeout: util.Timeout,
		}

		ginkgo.DeferCleanup(func() {
			client.CloseIdleConnections()
		})
	})

	ginkgo.Context("get current user endpoint", func() {
		ginkgo.It("should return user", func() {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/users/current", chatURL), nil)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", util.AdminUsername))

			resp, err := client.Do(req)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			ginkgo.DeferCleanup(func() {
				resp.Body.Close()
			})

			gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK))

			body, err := io.ReadAll(resp.Body)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			user := &commonhttp.UserDto{}
			err = json.Unmarshal(body, user)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			gomega.Expect(user).To(gomega.BeComparableTo(&commonhttp.UserDto{
				ID:       util.AdminID,
				Email:    util.AdminEmail,
				Username: util.AdminUsername,
			}))
		})
	})

	ginkgo.Context("get users endpoint", func() {
		ginkgo.It("should return user list", func() {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/users", chatURL), nil)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", util.AdminUsername))

			resp, err := client.Do(req)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			ginkgo.DeferCleanup(func() {
				resp.Body.Close()
			})

			gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK))

			body, err := io.ReadAll(resp.Body)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			page := commonhttp.Page[commonhttp.UserDto]{}
			err = json.Unmarshal(body, &page)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			gomega.Expect(page).To(gomega.BeComparableTo(commonhttp.Page[commonhttp.UserDto]{
				Items: []commonhttp.UserDto{
					{
						ID:       util.AdminID,
						Email:    util.AdminEmail,
						Username: util.AdminUsername,
					},
					{
						ID:       util.UserID,
						Email:    util.UserEmail,
						Username: util.UserUsername,
					},
				},
				Count: 2,
			}))
		})
	})
})
