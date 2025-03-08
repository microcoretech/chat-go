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

package api

import (
	"context"
	"testing"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"chat-go/test/integration/framework"
)

var (
	fwk *framework.Framework
)

func TestAPI(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "API Test Suite")
}

var _ = ginkgo.BeforeSuite(func() {
	ctx := context.Background()
	fwk = framework.NewFramework()
	gomega.Expect(fwk.Setup(ctx)).Should(gomega.Succeed())
})

var _ = ginkgo.AfterSuite(func() {
	ctx := context.Background()
	gomega.Expect(fwk.Teardown(ctx)).To(gomega.Succeed())
})
