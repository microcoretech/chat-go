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
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	projectDirEnv    = "PROJECT_DIR"
	migrationsDirEnv = "MIGRATIONS_DIR"
)

// GetProjectDir retrieves the project directory either from an environment variable or by searching
// for a Makefile.
func GetProjectDir() (string, error) {
	projectDir, found := os.LookupEnv(projectDirEnv)
	if found {
		return filepath.Dir(projectDir), nil
	}

	projectBaseDir, err := findMakefileDir()
	if err != nil {
		return "", err
	}

	return projectBaseDir, nil
}

// findMakefileDir traverses directories upward from the current directory until it finds a directory
// containing a Makefile.
func findMakefileDir() (string, error) {
	startDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("could not get current working directory: %w", err)
	}

	for {
		makefilePath := filepath.Join(startDir, "Makefile")
		if _, err := os.Stat(makefilePath); err == nil {
			return startDir, nil
		}

		parentDir := filepath.Dir(startDir)
		if parentDir == startDir {
			return "", errors.New("not able to locate Makefile")
		}

		startDir = parentDir
	}
}

// GetMigrationsDir retrieves the migrations directory either from an environment variable or using
// the absolute project path.
func GetMigrationsDir() (string, error) {
	migrationsDir, found := os.LookupEnv(migrationsDirEnv)
	if found {
		return filepath.Dir(migrationsDir), nil
	}

	projectDir, err := GetProjectDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(projectDir, "migrations"), nil
}
