package api

import (
	chaterrors "chat/internal/chat/errors"
	"chat/internal/common/common"
	"chat/internal/common/errors"
	"chat/internal/infrastructure/configs"
	"chat/internal/infrastructure/logger"
	usererrors "chat/internal/user/errors"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func ErrorHandler(log logger.Logger, environment configs.Environment) fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		var statusCode int

		switch err.(type) {
		case
			*fiber.Error:
			{
				fiberErr := err.(*fiber.Error)
				statusCode = fiberErr.Code
				switch statusCode {
				case http.StatusNotFound:
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
			ctx.Status(statusCode).JSON(errors.TruncateErrorData(errorData))
		} else {
			ctx.Status(statusCode).JSON(baseError)
		}

		return nil
	}
}
