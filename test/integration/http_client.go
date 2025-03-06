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

package integration

import (
	"net/http"
	"time"

	"chat-go/test/util"
)

var _ util.HTTPClient = (*TestHTTPClient)(nil)

type TestHTTPClient struct {
	env     *TestEnvironment
	timeout *time.Duration
}

func NewTestHTTPClient(env *TestEnvironment) *TestHTTPClient {
	return &TestHTTPClient{
		env: env,
	}
}

func (c *TestHTTPClient) WithTimeout(timeout time.Duration) *TestHTTPClient {
	c.timeout = &timeout
	return c
}

func (c *TestHTTPClient) Do(req *http.Request) (*http.Response, error) {
	msTimeout := -1
	if c.timeout != nil {
		msTimeout = int(c.timeout.Milliseconds())
	}
	return c.env.App().Test(req, msTimeout)
}
