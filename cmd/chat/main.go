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

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	chatdomain "chat-go/internal/chat/domain"
	chathttp "chat-go/internal/chat/http"
	chatrepository "chat-go/internal/chat/repository"
	chatwebsocket "chat-go/internal/chat/websocket"
	"chat-go/internal/common/repository"
	"chat-go/internal/infrastructure/api"
	"chat-go/internal/infrastructure/configs"
	"chat-go/internal/infrastructure/connector"
	"chat-go/internal/infrastructure/database/postgres"
	"chat-go/internal/infrastructure/logger/logrus"
	"chat-go/internal/infrastructure/validator"
	usercontract "chat-go/internal/user/contract"
	userdomain "chat-go/internal/user/domain"
	userhttp "chat-go/internal/user/http"
)

func main() {
	cfg, err := configs.NewConfig()
	if err != nil {
		panic(fmt.Errorf("errors on init config: %w", err))
	}

	log, err := logrus.NewLogger(cfg.LogLevel)
	if err != nil {
		panic(fmt.Errorf("error on init logger: %w", err))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbConn, err := postgres.NewPostgres(ctx, cfg.PostgresURI)
	if err != nil {
		log.Fatal(fmt.Errorf("error on connection to postgres: %w", err))
	}

	validate, err := validator.New()
	if err != nil {
		log.Fatalf("error on create validator: %s", err)
	}

	baseRepo := repository.NewBaseRepoImpl(dbConn)
	chatRepo := chatrepository.NewChatRepoImpl(dbConn)
	userChatRepo := chatrepository.NewUserChatRepoImpl(dbConn)
	messageRepo := chatrepository.NewMessageRepoImpl(dbConn)

	userService := userdomain.NewUserServiceImpl(cfg)
	userServiceContract := usercontract.NewUserServiceContractImpl(userService)
	chatService := chatdomain.NewChatServiceImpl(baseRepo, chatRepo, userChatRepo, userServiceContract)
	messageService := chatdomain.NewMessageServiceImpl(messageRepo, userServiceContract)

	eventHandler := chatwebsocket.NewEventHandler(validate, messageService)

	connector := connector.NewConnector(log, eventHandler)

	authMiddleware := userhttp.NewAuthMiddleware(userService)

	userController := userhttp.NewUserController(validate, authMiddleware, userService)
	chatController := chathttp.NewChatController(validate, authMiddleware, chatService, messageService, connector)

	server := api.NewHTTPServer(cfg, log, userController, chatController)

	ctx, cancel = signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		if err := connector.Start(ctx); err != nil {
			log.Errorf("Error on running connector: %s", err.Error())
			return err
		}

		log.Info("Connector gracefully stopped")

		return nil
	})

	eg.Go(func() error {
		if err := server.Start(ctx); err != nil {
			log.Errorf("Error on running server: %s", err.Error())
			return err
		}

		log.Info("Server gracefully stopped")

		return nil
	})

	if err = eg.Wait(); err != nil {
		log.Error(err)
	}
}
