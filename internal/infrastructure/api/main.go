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
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	fibercors "github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	fiberrecover "github.com/gofiber/fiber/v2/middleware/recover"

	"chat-go/internal/infrastructure/configs"
	"chat-go/internal/infrastructure/logger"
)

var (
	ServerAlreadyStartedError = errors.New("server already started")
)

type Server interface {
	Start(ctx context.Context) error
}

type httpServerImpl struct {
	log logger.Logger
	cfg *configs.Config

	app *fiber.App

	authMiddleware Middleware

	authController Controller
	userController Controller
	chatController Controller

	isStarted bool
}

func (s *httpServerImpl) init() {
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler(s.log, s.cfg.Environment),
	})

	app.Use(fiberlogger.New(fiberlogger.Config{
		TimeFormat: time.DateTime,
		Format:     "{\"status\":${status},\"latency\":\"${latency}\",\"method\":\"${method}\",\"url\":\"${url}\",\"ip\":\"${ip}\"}\n",
		Output:     s.log.Writer(),
	}))

	app.Use(fibercors.New())
	app.Use(fiberrecover.New())

	app.Get("/", rootHandler(s.cfg))
	app.Get("/healthz", healthzHandler)

	auth := app.Group("/auth")
	s.authController.SetupRoutes(auth)

	user := app.Group("/users", s.authMiddleware.Handler)
	s.userController.SetupRoutes(user)

	chat := app.Group("/chats", s.authMiddleware.Handler)
	s.chatController.SetupRoutes(chat)

	s.app = app
}

func (s *httpServerImpl) Start(ctx context.Context) error {
	if s.isStarted {
		return ServerAlreadyStartedError
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

func NewServer(
	cfg *configs.Config,
	log logger.Logger,
	authMiddleware Middleware,
	authController Controller,
	userController Controller,
	chatController Controller,
) Server {
	s := &httpServerImpl{
		cfg:            cfg,
		log:            log,
		authMiddleware: authMiddleware,
		authController: authController,
		userController: userController,
		chatController: chatController,
	}

	s.init()

	return s
}
