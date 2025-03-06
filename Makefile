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
ARTIFACTS ?= $(PROJECT_DIR)/bin
TOOLS_DIR := $(PROJECT_DIR)/hack/tools
BIN_DIR ?= $(PROJECT_DIR)/bin
ARTIFACTS_DIR ?= $(PROJECT_DIR)/artifacts

E2E_TARGET ?= $(PROJECT_DIR)/test/e2e/...
INTEGRATION_TARGET ?= $(PROJECT_DIR)/test/integration/...

BINARY ?= chat-go
IMAGE_NAME ?= chat-go

GIT_TAG ?= $(shell git describe --tags --dirty --always)

BINARY ?= chat-go
DOCKER_BUILDX_CMD ?= docker buildx
IMAGE_BUILD_CMD ?= $(DOCKER_BUILDX_CMD) build
IMAGE_BUILD_EXTRA_OPTS ?=
IMAGE_REGISTRY ?= mykhailobobrovskyi
IMAGE_NAME ?= chat-go
IMAGE_REPO ?= $(IMAGE_REGISTRY)/$(IMAGE_NAME)
IMAGE_TAG ?= $(IMAGE_REPO):$(GIT_TAG)
PLATFORMS ?= linux/amd64,linux/arm64

GO_CMD ?= go
GINKGO ?= $(BIN_DIR)/ginkgo

# Use go.mod go version as source.
GOLANGCI_LINT_VERSION ?= $(shell cd $(TOOLS_DIR); $(GO_CMD) list -m -f '{{.Version}}' github.com/golangci/golangci-lint)
GINKGO_VERSION ?= $(shell cd $(TOOLS_DIR); $(GO_CMD) list -m -f '{{.Version}}' github.com/onsi/ginkgo/v2)

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-24s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: verify
verify: gomod-verify ci-lint

.PHONY: gomod-verify
gomod-verify:
	$(GO_CMD) mod tidy
	git --no-pager diff --exit-code go.mod go.sum

.PHONY: ci-lint
ci-lint: golangci-lint
	$(GOLANGCI_LINT) run --timeout 15m0s

.PHONY: lint-fix
lint-fix: golangci-lint
	$(GOLANGCI_LINT) run --fix --timeout 15m0s

.PHONY: gomod-download
gomod-download:
	$(GO_CMD) mod download

##@ Tests

.PHONY: test-e2e
test-e2e: gomod-download ginkgo docker-build ## Run e2e tests.
	$(GINKGO) --race -v $(E2E_TARGET)

.PHONY: test-integration
test-integration: gomod-download ginkgo ## Run e2e tests.
	$(GINKGO) --race -v $(INTEGRATION_TARGET)

##@ Build

.PHONY: build
build:
	$(GO_CMD) build -o $(BIN_DIR)/$(BINARY) $(PROJECT_DIR)/cmd/chat/main.go

.PHONY: docker-build
docker-build:
	docker build -t $(IMAGE_NAME) .

##@ Tools
GOLANGCI_LINT = $(PROJECT_DIR)/bin/golangci-lint
.PHONY: golangci-lint
golangci-lint: ## Download golangci-lint locally if necessary.
	@GOBIN=$(PROJECT_DIR)/bin GO111MODULE=on $(GO_CMD) install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

.PHONY: ginkgo
ginkgo: ## Download ginkgo locally if necessary.
	@GOBIN=$(BIN_DIR) GO111MODULE=on $(GO_CMD) install github.com/onsi/ginkgo/v2/ginkgo@$(GINKGO_VERSION)
