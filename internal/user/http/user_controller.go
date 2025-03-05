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

package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"

	"chat-go/internal/common/domain"
	"chat-go/internal/common/errors"
	"chat-go/internal/common/http"
	"chat-go/internal/infrastructure/validator"
	"chat-go/internal/user/common"
	usererrors "chat-go/internal/user/errors"
)

type UserController struct {
	validate    validator.Validate
	userService UserService
}

func (c *UserController) getCurrentUser(ctx *fiber.Ctx) error {
	user, err := c.userService.GetCurrentUser(ctx.Context())
	if err != nil {
		return err
	}

	if user == nil {
		return usererrors.NewUserNotFoundError(nil)
	}

	return ctx.JSON(http.UserToDto(*user))
}

func (c *UserController) getUsers(ctx *fiber.Ctx) error {
	var query UserQuery

	if err := ctx.QueryParser(&query); err != nil {
		return errors.NewBadRequestError(common.UserDomain, err, nil)
	}

	if err := c.validate.Struct(common.UserDomain, &query); err != nil {
		return errors.NewValidationError(common.UserDomain, err, nil)
	}

	userFilter, err := UserFilterFromQuery(query)
	if err != nil {
		return err
	}

	users, count, err := c.userService.GetUsers(ctx.Context(), &userFilter)
	if err != nil {
		return err
	}

	return ctx.JSON(http.NewPage(
		lo.Map(users, func(user domain.User, _ int) http.UserDto {
			return http.UserToDto(user)
		}),
		count,
	))
}

func (c *UserController) getUser(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return errors.NewBadRequestError(common.UserDomain, err, map[string]any{"id": idStr})
	}

	users, _, err := c.userService.GetUsers(ctx.Context(), &domain.UserFilter{
		IDs: []uint64{id},
	})
	if err != nil {
		return err
	}

	if len(users) == 0 {
		return usererrors.NewUserNotFoundError(map[string]any{"id": id})
	}

	return ctx.JSON(http.UserToDto(users[0]))
}

func (c *UserController) SetupRoutes(r fiber.Router) {
	r.Get("/current", c.getCurrentUser)
	r.Get("", c.getUsers)
	r.Get("/:id", c.getUser)
}

func NewUserController(validate validator.Validate, userService UserService) *UserController {
	return &UserController{
		validate:    validate,
		userService: userService,
	}
}
