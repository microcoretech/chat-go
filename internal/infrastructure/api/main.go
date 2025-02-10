package api

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	fiberCors "github.com/gofiber/fiber/v2/middleware/cors"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	fiberRecover "github.com/gofiber/fiber/v2/middleware/recover"

	"mbobrovskyi/chat-go/internal/infrastructure/configs"
	"mbobrovskyi/chat-go/internal/infrastructure/logger"
)

var (
	ServerAlreadyStartedError = errors.New("server already started")
)

type Server interface {
	Start(ctx context.Context) error
}

type httpServerImpl struct {
	log     logger.Logger
	cfg     *configs.Config
	version string

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

	app.Use(fiberLogger.New(fiberLogger.Config{
		TimeFormat: time.DateTime,
		Format:     "[${time}] ${status} - ${latency} ${method} ${url}\n",
		Output:     log.Writer(),
	}))

	app.Use(fiberCors.New())

	app.Use(fiberRecover.New())

	app.Get("/", HealthHandler(s.version))

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
	version string,
	authMiddleware Middleware,
	authController Controller,
	userController Controller,
	chatController Controller,
) Server {
	s := &httpServerImpl{
		cfg:            cfg,
		log:            log,
		version:        version,
		authMiddleware: authMiddleware,
		authController: authController,
		userController: userController,
		chatController: chatController,
	}

	s.init()

	return s
}
