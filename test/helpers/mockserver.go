// Copyright MicroCore Tech
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

package helpers

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mockserver"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	MockServerImage = "mockserver/mockserver:5.15.0"
)

func RunMockserverContainer(ctx context.Context) (*mockserver.MockServerContainer, error) {
	mockserverSpecPath, err := GetMockserverSpecPath()
	if err != nil {
		return nil, err
	}

	return mockserver.Run(
		ctx,
		MockServerImage,
		testcontainers.WithWaitStrategy(wait.ForLog("INFO 1080 started on port: 1080")),
		testcontainers.WithEnv(map[string]string{
			"MOCKSERVER_INITIALIZATION_JSON_PATH": "/config/spec.json",
		}),
		testcontainers.WithHostConfigModifier(func(hostConfig *container.HostConfig) {
			hostConfig.Binds = append(hostConfig.Binds, fmt.Sprintf("%s:/config/spec.json", mockserverSpecPath))
		}),
	)
}

func BuildGetCurrentUserEndpoint(host, port string) string {
	return fmt.Sprintf("http://%s:%s/users/current", host, port)
}

func GetCurrentUserEndpointForContainer(ctx context.Context, mockserverContainer *mockserver.MockServerContainer) (string, error) {
	port, err := mockserverContainer.MappedPort(ctx, "1080/tcp")
	if err != nil {
		return "", err
	}

	host, err := mockserverContainer.Host(ctx)
	if err != nil {
		return "", err
	}

	return BuildGetCurrentUserEndpoint(host, port.Port()), nil
}

func BuildGetUsersEndpoint(host, port string) string {
	return fmt.Sprintf("http://%s:%s/users", host, port)
}

func GetUsersEndpointContainer(ctx context.Context, mockserverContainer *mockserver.MockServerContainer) (string, error) {
	port, err := mockserverContainer.MappedPort(ctx, "1080/tcp")
	if err != nil {
		return "", err
	}

	host, err := mockserverContainer.Host(ctx)
	if err != nil {
		return "", err
	}

	return BuildGetUsersEndpoint(host, port.Port()), nil
}
