# Copyright 2025 Mykhailo Bobrovskyi
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
BIN_DIR ?= $(PROJECT_DIR)/bin
TOOLS_DIR := $(PROJECT_DIR)/hack/tools

E2E_TARGET ?= $(PROJECT_DIR)/test/e2e/...
INTEGRATION_TARGET ?= $(PROJECT_DIR)/test/integration/...

GIT_TAG ?= $(shell git describe --tags --dirty --always)

BINARY ?= chat-go
IMAGE_REGISTRY ?= microcoretech
IMAGE_NAME ?= chat-go
IMAGE_REPO ?= $(IMAGE_REGISTRY)/$(IMAGE_NAME)
IMAGE_TAG ?= $(IMAGE_REPO):$(GIT_TAG)

GO ?= go
GOLANGCI_LINT ?= $(GO) tool github.com/golangci/golangci-lint/v2/cmd/golangci-lint
GINKGO ?= $(GO) tool github.com/onsi/ginkgo/v2/ginkgo

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-24s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: verify
verify: gomod-verify ci-lint

.PHONY: gomod-verify
gomod-verify:
	$(GO) mod tidy
	git --no-pager diff --exit-code go.mod go.sum

.PHONY: ci-lint
ci-lint:
	$(GOLANGCI_LINT) run --timeout 15m0s

.PHONY: lint-fix
lint-fix:
	$(GOLANGCI_LINT) run --fix --timeout 15m0s

.PHONY: gomod-download
gomod-download:
	$(GO) mod download

##@ Tests

.PHONY: test
test: test-integration test-e2e ## Run all tests.

.PHONY: test-integration
test-integration: gomod-download ## Run e2e tests.
	$(GINKGO) --race -v $(INTEGRATION_TARGET)

.PHONY: test-e2e
test-e2e: gomod-download docker-build ## Run e2e tests.
	IMAGE_TAG=$(IMAGE_TAG) $(GINKGO) -v $(E2E_TARGET)

##@ Build

.PHONY: build
build:
	$(GO) build -o $(BIN_DIR)/$(BINARY) $(PROJECT_DIR)/cmd/chat/main.go

.PHONY: docker-build
docker-build:
	docker build -t $(IMAGE_TAG) .
