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
	"chat-go/internal/infrastructure/database/redis"
	"chat-go/internal/infrastructure/logger/logrus"
	"chat-go/internal/infrastructure/validator"
	usercontract "chat-go/internal/user/contract"
	userdomain "chat-go/internal/user/domain"
	userhttp "chat-go/internal/user/http"
	userrepository "chat-go/internal/user/repository"
)

func main() {
	cfg, err := configs.NewConfig()
	if err != nil {
		panic(err)
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

	redisClient, err := redis.NewRedis(ctx, cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDb)
	if err != nil {
		log.Fatal(fmt.Errorf("error on connection to redis: %w", err))
	}

	validate, err := validator.New()
	if err != nil {
		log.Fatalf("error on create validator: %s", err)
	}

	baseRepo := repository.NewBaseRepoImpl(dbConn)
	userRepo := userrepository.NewUserRepoImpl(dbConn)
	userCredentialsRepo := userrepository.NewUserCredentialsRepoImpl(dbConn)
	sessionRepo := userrepository.NewSessionRepo(redisClient)
	chatRepo := chatrepository.NewChatRepoImpl(dbConn)
	userChatRepo := chatrepository.NewUserChatRepoImpl(dbConn)
	messageRepo := chatrepository.NewMessageRepoImpl(dbConn)

	authService := userdomain.NewAuthServiceImpl(baseRepo, userRepo, userCredentialsRepo, sessionRepo)
	userService := userdomain.NewUserServiceImpl(baseRepo, userRepo)
	userServiceContract := usercontract.NewUserServiceContractImpl(userService)
	chatService := chatdomain.NewChatServiceImpl(baseRepo, chatRepo, userChatRepo, userServiceContract)
	messageService := chatdomain.NewMessageServiceImpl(messageRepo, userServiceContract)

	eventHandler := chatwebsocket.NewEventHandler(validate, messageService)

	connector := connector.NewConnector(log, eventHandler)

	authMiddleware := userhttp.NewAuthMiddleware(authService)

	authController := userhttp.NewAuthController(validate, authService, authMiddleware)
	userController := userhttp.NewUserController(validate, userService)
	chatController := chathttp.NewChatController(validate, chatService, messageService, connector)

	server := api.NewServer(cfg, log, authMiddleware, authController, userController, chatController)

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
