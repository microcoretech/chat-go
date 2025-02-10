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
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"

	chaterrors "mbobrovskyi/chat-go/internal/chat/errors"
	"mbobrovskyi/chat-go/internal/common/common"
	"mbobrovskyi/chat-go/internal/common/errors"
	"mbobrovskyi/chat-go/internal/infrastructure/configs"
	"mbobrovskyi/chat-go/internal/infrastructure/logger"
	usererrors "mbobrovskyi/chat-go/internal/user/errors"
)

func ErrorHandler(log logger.Logger, environment configs.Environment) fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		var statusCode int

		switch errType := err.(type) {
		case
			*fiber.Error:
			{
				fiberErr := errType
				statusCode = fiberErr.Code
				if statusCode == http.StatusNotFound {
					err = errors.NewNotFoundError(common.CommonDomain)
				}
			}
		case
			*errors.BadRequestError,
			*errors.ValidationError,
			*chaterrors.IncorrectUsersCountError,
			*chaterrors.InvalidChatNameError,
			*chaterrors.InvalidChatTypeError:
			{
				statusCode = http.StatusBadRequest
			}

		case
			*errors.UnauthorizedError:
			{
				statusCode = http.StatusUnauthorized
			}

		case
			*errors.ForbiddenError:
			{
				statusCode = http.StatusForbidden
			}

		case
			*errors.NotFoundError,
			*usererrors.UserNotFoundError:
			{
				statusCode = http.StatusNotFound
			}

		case
			*usererrors.UserAlreadyCreatedError:
			{
				statusCode = http.StatusConflict
			}

		case
			*usererrors.UserNotCreatedError,
			*errors.DatabaseError:
			{
				statusCode = http.StatusInternalServerError
			}

		default:
			{
				statusCode = http.StatusInternalServerError
				err = errors.NewUndefinedError(err)
			}
		}

		baseError := err.(errors.BaseError)

		errorData := baseError.GetErrorData()
		errorData.Data["path"] = fmt.Sprintf("%s %s", ctx.Method(), ctx.Request().URI().String())

		jsonErr, _ := json.Marshal(err)
		if statusCode < 500 {
			log.Debug(string(jsonErr))
		} else {
			log.Error(string(jsonErr))
		}

		if environment != configs.DevelopmentEnvironment {
			_ = ctx.Status(statusCode).JSON(errors.TruncateErrorData(errorData))
		} else {
			_ = ctx.Status(statusCode).JSON(baseError)
		}

		return nil
	}
}
