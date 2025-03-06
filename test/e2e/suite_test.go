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
	"context"
	"fmt"
	"testing"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"chat-go/test/util"
)

var (
	infra         *util.Infrastructure
	chatContainer testcontainers.Container
	chatURL       string
)

func TestAPI(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "API Test Suite")
}

var _ = ginkgo.BeforeSuite(func() {
	ctx := context.Background()

	infra = util.NewInfrastructure()
	gomega.Expect(infra.Init(ctx)).Should(gomega.Succeed())

	var err error

	postgresContainerIP, err := infra.PostgresContainer().ContainerIP(ctx)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	mockserverContainerIP, err := infra.MockserverContainer().ContainerIP(ctx)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	chatContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        util.ChatImage,
			ExposedPorts: []string{"8080/tcp"},
			WaitingFor:   wait.ForHTTP("/").WithPort("8080/tcp"),
			Env: map[string]string{
				"POSTGRES_URI":              util.BuildPostgresURI(postgresContainerIP, "5432"),
				"GET_CURRENT_USER_ENDPOINT": util.BuildGetCurrentUserEndpoint(mockserverContainerIP, "1080"),
				"GET_USERS_ENDPOINT":        util.BuildGetUsersEndpoint(mockserverContainerIP, "1080"),
			},
		},
		Started: true,
	})
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	chatPort, err := chatContainer.MappedPort(ctx, "8080/tcp")
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	chatHost, err := chatContainer.Host(ctx)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	chatURL = fmt.Sprintf("http://%s:%s", chatHost, chatPort.Port())
})

var _ = ginkgo.AfterSuite(func() {
	gomega.Expect(testcontainers.TerminateContainer(chatContainer)).To(gomega.Succeed())
	gomega.Expect(infra.Cleanup(context.Background())).To(gomega.Succeed())
})
