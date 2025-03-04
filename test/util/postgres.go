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

package util

import (
	"context"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	migratepostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	testcontainerspostgres "github.com/testcontainers/testcontainers-go/modules/postgres"

	"chat-go/internal/infrastructure/database/postgres"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	PostgresImage = "postgres:17-alpine"

	DbName = "test-db"
	DbUser = "user"
	DbPass = "password"
)

func PostgresURL(ctx context.Context, postgresContainer *testcontainerspostgres.PostgresContainer) (string, error) {
	chatPort, err := postgresContainer.MappedPort(ctx, "5432/tcp")
	if err != nil {
		return "", err
	}

	chatHost, err := postgresContainer.Host(ctx)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		DbUser, DbPass, chatHost, chatPort.Port(), DbName,
	), nil
}

func Migrate(ctx context.Context, postgresContainer *testcontainerspostgres.PostgresContainer) error {
	postgresURL, err := PostgresURL(ctx, postgresContainer)
	if err != nil {
		return err
	}

	db, err := postgres.NewPostgres(ctx, postgresURL)
	if err != nil {
		return err
	}
	defer db.Close()

	driver, err := migratepostgres.WithInstance(db, &migratepostgres.Config{})
	if err != nil {
		return err
	}

	migrationsDir, err := GetMigrationsDir()
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", migrationsDir), DbName, driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil {
		return err
	}

	return nil
}
