package http

import (
	"chat/internal/common/domain"
	"chat/internal/common/errors"
	"chat/internal/common/http"
	"chat/internal/infrastructure/validator"
	"chat/internal/user/common"
	usererrors "chat/internal/user/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"strconv"
)

type UserController struct {
	validate    validator.Validate
	userService UserService
}

func (c *UserController) getCurrentUser(ctx *fiber.Ctx) error {
	session := ctx.Context().UserValue("session").(*domain.Session)

	user, err := c.userService.GetUser(ctx.Context(), session.User.ID)
	if err != nil {
		return err
	}

	if user == nil {
		return usererrors.NewUserNotFoundError(map[string]any{"id": session.User.ID})
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

	user, err := c.userService.GetUser(ctx.Context(), id)
	if err != nil {
		return err
	}

	if user == nil {
		return usererrors.NewUserNotFoundError(map[string]any{"id": id})
	}

	return ctx.JSON(http.UserToDto(*user))
}

func (c *UserController) SetupRoutes(r fiber.Router) {
	r.Get("/current", c.getCurrentUser)
	r.Get("", c.getUsers)
	r.Get("/:id", c.getUser)
}

func NewUserController(
	validate validator.Validate,
	userService UserService,
) *UserController {
	return &UserController{
		validate:    validate,
		userService: userService,
	}
}
