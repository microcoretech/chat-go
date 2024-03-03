package main

import (
	chatdomain "chat/internal/chat/domain"
	chathttp "chat/internal/chat/http"
	chatrepository "chat/internal/chat/repository"
	chatwebsocket "chat/internal/chat/websocket"
	"chat/internal/common/repository"
	"chat/internal/infrastructure/api"
	"chat/internal/infrastructure/configs"
	"chat/internal/infrastructure/connector"
	"chat/internal/infrastructure/database/postgres"
	"chat/internal/infrastructure/database/redis"
	"chat/internal/infrastructure/logger/logrus"
	"chat/internal/infrastructure/validator"
	usercontract "chat/internal/user/contract"
	userdomain "chat/internal/user/domain"
	userhttp "chat/internal/user/http"
	userrepository "chat/internal/user/repository"
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := configs.NewConfig()
	if err != nil {
		panic(err)
	}

	fileVersion, err := os.ReadFile("VERSION")
	if err != nil {
		panic(fmt.Errorf("error on reading VERSION file: %w", err))
	}

	version := string(fileVersion)

	log, err := logrus.NewLogger(cfg.LogLevel)
	if err != nil {
		panic(fmt.Errorf("error on init logger: %w", err))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbConn, err := postgres.NewPostgres(ctx, cfg.PostgresUri)
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

	server := api.NewServer(
		cfg,
		log,
		version,
		authMiddleware,
		authController,
		userController,
		chatController,
	)

	ctx, cancel = signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

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
