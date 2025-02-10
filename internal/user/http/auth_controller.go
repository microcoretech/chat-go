package http

import (
	"github.com/gofiber/fiber/v2"

	"mbobrovskyi/chat-go/internal/common/domain"
	"mbobrovskyi/chat-go/internal/common/errors"
	"mbobrovskyi/chat-go/internal/infrastructure/api"
	"mbobrovskyi/chat-go/internal/infrastructure/validator"
	"mbobrovskyi/chat-go/internal/user/common"
)

type AuthController struct {
	validate       validator.Validate
	authService    AuthService
	authMiddleware api.Middleware
}

func (c *AuthController) signIn(ctx *fiber.Ctx) error {
	var signInRequest SignInRequest

	if err := ctx.BodyParser(&signInRequest); err != nil {
		return errors.NewBadRequestError(common.UserDomain, err, nil)
	}

	if err := c.validate.Struct(common.UserDomain, &signInRequest); err != nil {
		return err
	}

	token, err := c.authService.SignIn(ctx.Context(), signInRequest.Username, signInRequest.Password)
	if err != nil {
		return err
	}

	return ctx.JSON(TokenToDto(*token))
}

func (c *AuthController) signUp(ctx *fiber.Ctx) error {
	var signUpRequest SignUpRequest

	if err := ctx.BodyParser(&signUpRequest); err != nil {
		return errors.NewBadRequestError(common.UserDomain, err, nil)
	}

	if err := c.validate.Struct(common.UserDomain, &signUpRequest); err != nil {
		return err
	}

	token, err := c.authService.SignUp(ctx.Context(), UserFromSignUpRequest(signUpRequest), signUpRequest.Password)
	if err != nil {
		return err
	}

	return ctx.JSON(TokenToDto(*token))
}

func (c *AuthController) signOut(ctx *fiber.Ctx) error {
	session := ctx.Context().UserValue("session").(*domain.Session)
	token := ctx.Context().UserValue("token").(string)
	if err := c.authService.SignOut(ctx.Context(), session.User.ID, token); err != nil {
		return err
	}
	return nil
}

func (c *AuthController) SetupRoutes(r fiber.Router) {
	r.Post("sign-in", c.signIn)
	r.Post("sign-up", c.signUp)
	r.Post("sign-out", c.authMiddleware.Handler, c.signOut)
}

func NewAuthController(
	validate validator.Validate,
	authService AuthService,
	authMiddleware api.Middleware,
) *AuthController {
	return &AuthController{
		validate:       validate,
		authService:    authService,
		authMiddleware: authMiddleware,
	}
}
