// Copyright MicroCore Tech
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
	"time"

	"github.com/gofiber/fiber/v2"
	fibercors "github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	fiberrecover "github.com/gofiber/fiber/v2/middleware/recover"

	"chat-go/internal/infrastructure/configs"
	"chat-go/internal/infrastructure/logger"
)

func NewApp(cfg *configs.Config, log logger.Logger, controllers ...Controller) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler(log, cfg.Environment),
	})

	app.Use(fiberlogger.New(fiberlogger.Config{
		TimeFormat: time.DateTime,
		Format:     "{\"status\":${status},\"latency\":\"${latency}\",\"method\":\"${method}\",\"url\":\"${url}\",\"ip\":\"${ip}\"}\n",
		Output:     log.Writer(),
	}))

	app.Use(fibercors.New())
	app.Use(fiberrecover.New())

	app.Get("/", rootHandler(cfg))
	app.Get("/healthz", healthzHandler)

	for _, controller := range controllers {
		controller.SetupRoutes(app)
	}

	return app
}
