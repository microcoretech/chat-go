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

package framework

import (
	"context"
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/onsi/ginkgo/v2/dsl/core"
	"github.com/sirupsen/logrus"

	chatdomain "chat-go/internal/chat/domain"
	chathttp "chat-go/internal/chat/http"
	chatrepository "chat-go/internal/chat/repository"
	chatwebsocket "chat-go/internal/chat/websocket"
	"chat-go/internal/common/repository"
	"chat-go/internal/infrastructure/api"
	"chat-go/internal/infrastructure/configs"
	"chat-go/internal/infrastructure/connector"
	"chat-go/internal/infrastructure/database/postgres"
	"chat-go/internal/infrastructure/logger"
	loggerlogrus "chat-go/internal/infrastructure/logger/logrus"
	"chat-go/internal/infrastructure/validator"
	usercontract "chat-go/internal/user/contract"
	userdomain "chat-go/internal/user/domain"
	userhttp "chat-go/internal/user/http"
	"chat-go/test/util"
)

type Framework struct {
	*util.Infrastructure

	cfg      *configs.Config
	log      *logrus.Logger
	dbConn   *sql.DB
	validate validator.Validate

	baseRepo     *repository.BaseRepoImpl
	userChatRepo *chatrepository.UserChatRepoImpl
	chatRepo     *chatrepository.ChatRepoImpl

	userService    *userdomain.UserServiceImpl
	chatService    *chatdomain.ChatServiceImpl
	messageService *chatdomain.MessageServiceImpl

	userServiceContract *usercontract.UserServiceContractImpl

	authMiddleware *userhttp.AuthMiddleware

	connector      *connector.ConnectorImpl
	userController *userhttp.UserController
	chatController *chathttp.ChatController
	eventHandler   *chatwebsocket.EventHandler

	app *fiber.App
}

func NewFramework() *Framework {
	return &Framework{
		Infrastructure: util.NewInfrastructure(),
	}
}

func (f *Framework) Setup(ctx context.Context) error {
	err := f.Infrastructure.Setup(ctx)
	if err != nil {
		return err
	}

	f.cfg = &configs.Config{
		LogLevel: logger.DebugLevel,
	}

	f.log, err = loggerlogrus.NewLogger(f.cfg.LogLevel)
	if err != nil {
		return err
	}
	f.log.SetOutput(core.GinkgoWriter)

	f.cfg.PostgresURI, err = util.PostgresURIForContainer(ctx, f.PostgresContainer())
	if err != nil {
		return err
	}

	f.cfg.GetCurrentUserEndpoint, err = util.GetCurrentUserEndpointForContainer(ctx, f.MockserverContainer())
	if err != nil {
		return err
	}

	f.cfg.GetUsersEndpoint, err = util.GetUsersEndpointContainer(ctx, f.MockserverContainer())
	if err != nil {
		return err
	}

	f.dbConn, err = postgres.NewPostgres(ctx, f.cfg.PostgresURI)
	if err != nil {
		return err
	}

	f.validate, err = validator.New()
	if err != nil {
		return err
	}

	f.baseRepo = repository.NewBaseRepoImpl(f.dbConn)
	f.chatRepo = chatrepository.NewChatRepoImpl(f.dbConn)
	f.userChatRepo = chatrepository.NewUserChatRepoImpl(f.dbConn)
	f.userService = userdomain.NewUserServiceImpl(f.cfg)
	f.chatService = chatdomain.NewChatServiceImpl(f.baseRepo, f.chatRepo, f.userChatRepo, f.userService)
	f.userServiceContract = usercontract.NewUserServiceContractImpl(f.userService)
	f.authMiddleware = userhttp.NewAuthMiddleware(f.userService)
	f.connector = connector.NewConnector(f.log, f.eventHandler)
	f.userController = userhttp.NewUserController(f.validate, f.authMiddleware, f.userService)
	f.chatController = chathttp.NewChatController(f.validate, f.authMiddleware, f.chatService, f.messageService, f.connector)
	f.app = api.NewApp(f.cfg, f.log, f.userController, f.chatController)
	f.eventHandler = chatwebsocket.NewEventHandler(f.validate, f.messageService)

	return nil
}

func (f *Framework) App() *fiber.App {
	return f.app
}

func (f *Framework) GetChatController() *chathttp.ChatController {
	return f.chatController
}
