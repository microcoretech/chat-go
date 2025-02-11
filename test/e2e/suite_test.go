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
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	chatImage = "chat-go"

	postgresImage = "postgres:17-alpine"
	redisImage    = "redis:7-alpine"

	dbName = "test-db"
	dbUser = "user"
	dbPass = "password"
)

var (
	postgresContainer *postgres.PostgresContainer
	redisContainer    *redis.RedisContainer
	chatContainer     testcontainers.Container

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
		postgresImage,
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPass),
		testcontainers.WithWaitStrategy(wait.ForListeningPort("5432/tcp")),
	)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	redisContainer, err = redis.Run(
		ctx,
		redisImage,
		testcontainers.WithWaitStrategy(wait.ForListeningPort("6379/tcp")),
	)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	postgresContainerIP, err := postgresContainer.ContainerIP(ctx)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	redisContainerIP, err := redisContainer.ContainerIP(ctx)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	chatContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        chatImage,
			ExposedPorts: []string{"8080/tcp"},
			WaitingFor:   wait.ForHTTP("/").WithPort("8080/tcp"),
			Env: map[string]string{
				"POSTGRES_URI": fmt.Sprintf(
					"postgresql://%s:%s@%s:5432/%s?sslmode=disable",
					dbUser, dbPass, postgresContainerIP, dbName,
				),
				"REDIS_ADDR": fmt.Sprintf("%s:6379", redisContainerIP),
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
	gomega.Expect(testcontainers.TerminateContainer(postgresContainer)).To(gomega.Succeed())
	gomega.Expect(testcontainers.TerminateContainer(redisContainer)).To(gomega.Succeed())
})
