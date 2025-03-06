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
	"context"
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/onsi/ginkgo/v2/dsl/core"
	"github.com/sirupsen/logrus"

	chatdomain "chat-go/internal/chat/domain"
	chathttp "chat-go/internal/chat/http"
	chatrepository "chat-go/internal/chat/repository"
	"chat-go/internal/common/repository"
	"chat-go/internal/infrastructure/api"
	"chat-go/internal/infrastructure/configs"
	"chat-go/internal/infrastructure/connector"
	"chat-go/internal/infrastructure/database/postgres"
	"chat-go/internal/infrastructure/validator"
	usercontract "chat-go/internal/user/contract"
	userdomain "chat-go/internal/user/domain"
	userhttp "chat-go/internal/user/http"
	"chat-go/test/util"
)

type TestEnvironment struct {
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

	userController *userhttp.UserController
	chatController *chathttp.ChatController

	connector *connector.ConnectorImpl

	httpServer *api.HTTPServer
}

func (e *TestEnvironment) Init(ctx context.Context) error {
	err := e.Infrastructure.Init(ctx)
	if err != nil {
		return err
	}

	e.cfg = &configs.Config{
		LogLevel:               configs.DevelopmentEnvironment.String(),
		GetCurrentUserEndpoint: "",
		GetUsersEndpoint:       "",
	}

	e.log = logrus.New()
	e.log.SetOutput(core.GinkgoWriter)

	e.cfg.PostgresURI, err = util.PostgresURIForContainer(ctx, e.Infrastructure.PostgresContainer())
	if err != nil {
		return err
	}

	e.cfg.GetCurrentUserEndpoint, err = util.GetCurrentUserEndpointForContainer(ctx, e.MockserverContainer())
	if err != nil {
		return err
	}

	e.cfg.GetUsersEndpoint, err = util.GetUsersEndpointContainer(ctx, e.MockserverContainer())
	if err != nil {
		return err
	}

	e.dbConn, err = postgres.NewPostgres(ctx, e.cfg.PostgresURI)
	if err != nil {
		return err
	}

	e.validate, err = validator.New()
	if err != nil {
		return err
	}

	e.baseRepo = repository.NewBaseRepoImpl(e.dbConn)
	e.chatRepo = chatrepository.NewChatRepoImpl(e.dbConn)
	e.userChatRepo = chatrepository.NewUserChatRepoImpl(e.dbConn)
	e.userService = userdomain.NewUserServiceImpl(e.cfg)
	e.chatService = chatdomain.NewChatServiceImpl(e.baseRepo, e.chatRepo, e.userChatRepo, e.userService)
	e.userServiceContract = usercontract.NewUserServiceContractImpl(e.userService)
	e.authMiddleware = userhttp.NewAuthMiddleware(e.userService)
	e.userController = userhttp.NewUserController(e.validate, e.authMiddleware, e.userService)
	e.chatController = chathttp.NewChatController(e.validate, e.authMiddleware, e.chatService, e.messageService, e.connector)
	e.httpServer = api.NewHTTPServer(e.cfg, e.log, e.userController, e.chatController)

	return nil
}

func (e *TestEnvironment) App() *fiber.App {
	return e.httpServer.App()
}

func (e *TestEnvironment) GetChatController() *chathttp.ChatController {
	return e.chatController
}

func NewTestEnvironment() *TestEnvironment {
	return &TestEnvironment{
		Infrastructure: util.NewInfrastructure(),
	}
}
