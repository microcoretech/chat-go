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

package helpers

import (
	"context"

	"github.com/testcontainers/testcontainers-go/modules/mockserver"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type Infrastructure struct {
	postgresContainer   *postgres.PostgresContainer
	mockserverContainer *mockserver.MockServerContainer
}

func NewInfrastructure() *Infrastructure {
	return &Infrastructure{}
}

func (i *Infrastructure) Setup(ctx context.Context) error {
	var err error

	i.postgresContainer, err = RunPostgresContainer(ctx)
	if err != nil {
		return err
	}

	i.mockserverContainer, err = RunMockserverContainer(ctx)
	if err != nil {
		return err
	}

	err = Migrate(ctx, i.postgresContainer)
	if err != nil {
		return err
	}

	return nil
}

func (i *Infrastructure) Teardown(ctx context.Context) error {
	if i.mockserverContainer != nil {
		if err := i.mockserverContainer.Terminate(ctx); err != nil {
			return err
		}
	}

	if i.postgresContainer != nil {
		if err := i.postgresContainer.Terminate(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (i *Infrastructure) PostgresContainer() *postgres.PostgresContainer {
	return i.postgresContainer
}

func (i *Infrastructure) MockserverContainer() *mockserver.MockServerContainer {
	return i.mockserverContainer
}
