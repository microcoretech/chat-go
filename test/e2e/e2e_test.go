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
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"chat-go/internal/common/common"
	"chat-go/internal/infrastructure/api"
	"chat-go/test/util"
)

var _ = ginkgo.Describe("Chat", func() {
	var client *http.Client

	ginkgo.BeforeEach(func() {
		client = &http.Client{
			Timeout: 5 * time.Second,
		}
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
			gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK))

			body, err := io.ReadAll(resp.Body)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(string(body)).To(gomega.Equal("OK"))
		})
	})
})
