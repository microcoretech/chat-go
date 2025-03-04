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

	"github.com/docker/docker/api/types/container"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mockserver"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"

	"chat-go/test/util"
)

var (
	postgresContainer   *postgres.PostgresContainer
	redisContainer      *redis.RedisContainer
	mockserverContainer *mockserver.MockServerContainer
	chatContainer       testcontainers.Container

	chatURL string
)

func TestAPI(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "API Test Suite")
}

var _ = ginkgo.BeforeSuite(func() {
	ctx := context.Background()

	var err error

	postgresContainer, err = postgres.Run(
		ctx,
		util.PostgresImage,
		postgres.WithDatabase(util.DbName),
		postgres.WithUsername(util.DbUser),
		postgres.WithPassword(util.DbPass),
		testcontainers.WithWaitStrategy(wait.ForListeningPort("5432/tcp")),
	)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	redisContainer, err = redis.Run(
		ctx,
		util.RedisImage,
		testcontainers.WithWaitStrategy(wait.ForListeningPort("6379/tcp")),
	)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	mockserverSpecPath, err := util.GetMockserverSpecPath()
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	mockserverContainer, err = mockserver.Run(
		ctx,
		util.MockServerImage,
		testcontainers.WithWaitStrategy(wait.ForLog("INFO 1080 started on port: 1080")),
		testcontainers.WithEnv(map[string]string{
			"MOCKSERVER_INITIALIZATION_JSON_PATH": "/config/spec.json",
		}),
		testcontainers.WithHostConfigModifier(func(hostConfig *container.HostConfig) {
			hostConfig.Binds = append(hostConfig.Binds, fmt.Sprintf("%s:/config/spec.json", mockserverSpecPath))
		}),
	)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	postgresContainerIP, err := postgresContainer.ContainerIP(ctx)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	err = util.Migrate(ctx, postgresContainer)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	redisContainerIP, err := redisContainer.ContainerIP(ctx)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	mockserverContainerIP, err := mockserverContainer.ContainerIP(ctx)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	chatContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        util.ChatImage,
			ExposedPorts: []string{"8080/tcp"},
			WaitingFor:   wait.ForHTTP("/").WithPort("8080/tcp"),
			Env: map[string]string{
				"POSTGRES_URI": fmt.Sprintf(
					"postgresql://%s:%s@%s:5432/%s?sslmode=disable",
					util.DbUser, util.DbPass, postgresContainerIP, util.DbName,
				),
				"REDIS_ADDR":                fmt.Sprintf("%s:6379", redisContainerIP),
				"GET_CURRENT_USER_ENDPOINT": fmt.Sprintf("http://%s:1080/users/current", mockserverContainerIP),
				"GET_USERS_ENDPOINT":        fmt.Sprintf("http://%s:1080/users", mockserverContainerIP),
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
	gomega.Expect(testcontainers.TerminateContainer(mockserverContainer)).To(gomega.Succeed())
	gomega.Expect(testcontainers.TerminateContainer(postgresContainer)).To(gomega.Succeed())
	gomega.Expect(testcontainers.TerminateContainer(redisContainer)).To(gomega.Succeed())
})
