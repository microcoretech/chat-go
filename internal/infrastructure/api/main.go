// Copyright 2025 MicroCore Tech
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
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"

	"chat-go/internal/infrastructure/configs"
	"chat-go/internal/infrastructure/logger"
)

var (
	ErrServerAlreadyStarted = errors.New("server already started")
)

type HTTPServer struct {
	cfg *configs.Config
	app *fiber.App

	isStarted bool
}

func (s *HTTPServer) App() *fiber.App {
	return s.app
}

func (s *HTTPServer) Start(ctx context.Context) error {
	if s.isStarted {
		return ErrServerAlreadyStarted
	}

	s.isStarted = true
	defer func() {
		s.isStarted = false
	}()

	go func() {
		<-ctx.Done()
		_ = s.app.ShutdownWithTimeout(time.Minute)
	}()

	if err := s.app.Listen(s.cfg.HTTPServerAddr); err != nil {
		return err
	}

	return nil
}

func NewHTTPServer(cfg *configs.Config, log logger.Logger, controllers ...Controller) *HTTPServer {
	return &HTTPServer{
		app: NewApp(cfg, log, controllers...),
		cfg: cfg,
	}
}
